package service

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"app/internal/model/ws"
	"app/internal/repository"

	"github.com/gofiber/contrib/v3/websocket"
)

const (
	writeTimeout  = 5 * time.Second
	readTimeout   = writeTimeout
	messageBuffer = 64
)

type client struct {
	userID int
	conn   *websocket.Conn
	send   chan []byte
	ctx    context.Context
}

func NewClient(userID int, conn *websocket.Conn, ctx context.Context) *client {
	return &client{
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, messageBuffer), // buffer for slow connections
		ctx:    ctx,
	}
}

func (c *client) WritePump(hub HubService, pingInterval time.Duration) {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *client) ReadPump(hub HubService, pongTimeout time.Duration) {
	defer func() {
		hub.Unregister(c)
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongTimeout))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var event ws.Event
		if err = json.Unmarshal(message, &event); err != nil {
			continue
		}

		switch event.Type {
		case ws.EventTypingStart:
			var payload ws.TypingPayloadInbound
			if err := json.Unmarshal(event.Payload, &payload); err != nil {
				continue
			}
			hub.HandleTypingStart(c.ctx, c, payload)
		case ws.EventTypingStop:
			var payload ws.TypingPayloadInbound
			if err := json.Unmarshal(event.Payload, &payload); err != nil {
				continue
			}
			hub.HandleTypingStop(c.ctx, c, payload)
		case ws.EventMessageRead:
			var payload ws.MessageReadPayload
			if err := json.Unmarshal(event.Payload, &payload); err != nil {
				continue
			}
			hub.HandleMessageRead(c.ctx, c, payload)
		}
	}
}

type hub struct {
	clients          map[int]map[*client]bool
	register         chan *client
	unregister       chan *client
	mu               sync.RWMutex
	messageService   MessageService
	conversationRepo repository.ConversationRepository
	userRepo         repository.UserRepository
	stop             chan struct{}
}

type HubService interface {
	Run()
	Register(client *client)
	Unregister(client *client)
	IsOnline(userID int) bool
	BroadcastToUser(ref string, userID int, payload any)
	BroadcastToConversation(ctx context.Context, userID int, ref string, conversationID int, payload any)
	HandleTypingStart(ctx context.Context, client *client, payload ws.TypingPayloadInbound)
	HandleTypingStop(ctx context.Context, client *client, payload ws.TypingPayloadInbound)
	HandleMessageRead(ctx context.Context, client *client, payload ws.MessageReadPayload)
	Stop()
}

func NewHubService(conversationRepo repository.ConversationRepository, userRepo repository.UserRepository, messageService MessageService) HubService {
	return &hub{
		clients:          make(map[int]map[*client]bool),
		register:         make(chan *client),
		unregister:       make(chan *client),
		mu:               sync.RWMutex{},
		messageService:   messageService,
		conversationRepo: conversationRepo,
		userRepo:         userRepo,
		stop:             make(chan struct{}),
	}
}

func (h *hub) Stop() {
	h.stop <- struct{}{}
}

func (h *hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			if h.clients[conn.userID] == nil {
				h.clients[conn.userID] = make(map[*client]bool)
			}
			h.clients[conn.userID][conn] = true
			h.mu.Unlock()
			h.broadcastOnlineStatus(conn.ctx, conn.userID, true)

		case conn := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[conn.userID]; ok {
				delete(clients, conn)
				close(conn.send)
				if len(clients) == 0 {
					delete(h.clients, conn.userID)
				}
			}
			h.mu.Unlock()
			h.broadcastOnlineStatus(conn.ctx, conn.userID, false)

		case <-h.stop:
			h.mu.Lock()
			for _, clients := range h.clients {
				for client := range clients {
					close(client.send)
				}
			}
			h.mu.Unlock()
			close(h.register)
			return
		}
	}
}

func (h *hub) Register(client *client) {
	h.register <- client
}

func (h *hub) Unregister(client *client) {
	h.unregister <- client
}

func (h *hub) IsOnline(userID int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

func (h *hub) BroadcastToUser(ref string, userID int, payload any) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		broadcastError(h, userID, ref, err)
		return
	}

	event := ws.Event{
		Type:    ref,
		Payload: payloadBytes,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		broadcastError(h, userID, ref, err)
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	conns := h.clients[userID]

	for client := range conns {
		select {
		case client.send <- eventBytes:
		default:
			// client connection is too slow so disconnect client
			// client should reconnect again
			h.Unregister(client)
		}
	}
}

func (h *hub) BroadcastToConversation(ctx context.Context, userID int, ref string, conversationID int, payload any) {
	members, err := h.conversationRepo.GetMembers(ctx, conversationID)
	if err != nil {
		broadcastError(h, userID, ref, err)
		return
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		broadcastError(h, userID, ref, err)
		return
	}

	event := ws.Event{
		Type:    ref,
		Payload: payloadBytes,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		broadcastError(h, userID, ref, err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, member := range members {
		conns := h.clients[member.UserID]

		for client := range conns {
			select {
			case client.send <- eventBytes:
			default:
				// client connection is too slow so disconnect client
				// client should reconnect again
				h.Unregister(client)
			}
		}
	}
}

func (h *hub) HandleTypingStart(ctx context.Context, client *client, payload ws.TypingPayloadInbound) {
	payloadResp := ws.TypingPayloadOutbound{
		ConversationID: payload.ConversationID,
		UserID:         client.userID,
	}

	h.BroadcastToConversation(ctx, client.userID, ws.EventUserTypingStart, payload.ConversationID, payloadResp)
}

func (h *hub) HandleTypingStop(ctx context.Context, client *client, payload ws.TypingPayloadInbound) {
	payloadResp := ws.TypingPayloadOutbound{
		ConversationID: payload.ConversationID,
		UserID:         client.userID,
	}

	h.BroadcastToConversation(ctx, client.userID, ws.EventUserTypingStop, payload.ConversationID, payloadResp)
}

func (h *hub) HandleMessageRead(ctx context.Context, client *client, payload ws.MessageReadPayload) {
	err := h.messageService.MarkAsRead(ctx, client.userID, payload.ConversationID, payload.MessageID)

	if err != nil {
		broadcastError(h, client.userID, ws.EventMessageRead, err)
		return
	}

	payloadResp := ws.MessageSeenPayload{
		UserID:         client.userID,
		MessageID:      payload.MessageID,
		ConversationID: payload.ConversationID,
	}

	h.BroadcastToConversation(ctx, client.userID, ws.EventMessageSeen, payloadResp.ConversationID, payloadResp)
}

func broadcastError(h HubService, userID int, ref string, err error) {
	payload := ws.ErrorPayload{Ref: ref}

	if serviceError, ok := errors.AsType[*Error](err); ok {
		payload.Message = serviceError.Message
	} else {
		payload.Message = "something went wrong"
	}

	h.BroadcastToUser(ws.EventError, userID, payload)
}

func (h *hub) broadcastOnlineStatus(ctx context.Context, userID int, online bool) {
	conversations, err := h.conversationRepo.GetByUserID(ctx, userID)

	if err != nil {
		broadcastError(h, userID, ws.EventBroadcastUserStatus, err)
		return
	}

	var payload any
	if online {
		payload = ws.UserOnlinePayload{
			UserID: userID,
		}
	} else {
		user, err := h.userRepo.UpdateLastSeenAt(ctx, userID)

		if err != nil {
			return
		}

		payload = ws.UserOfflinePayload{
			UserID:     userID,
			LastSeenAt: *user.LastSeenAt,
		}
	}

	conversationIDs := make([]int, len(conversations))

	for i, conversation := range conversations {
		conversationIDs[i] = conversation.ID
	}

	members, err := h.conversationRepo.GetMembersByConversationIDs(ctx, conversationIDs)

	if err != nil {
		return
	}

	memberMap := make(map[int]bool)

	for _, member := range members {
		memberMap[member.UserID] = true
	}

	// verify user status before broadcasting
	if (!online && h.IsOnline(userID)) || (online && !h.IsOnline(userID)) {
		return
	}

	for memberID := range memberMap {
		if memberID != userID {
			if online {
				h.BroadcastToUser(ws.EventUserOnline, memberID, payload)
			} else {
				h.BroadcastToUser(ws.EventUserOffline, memberID, payload)
			}
		}
	}
}
