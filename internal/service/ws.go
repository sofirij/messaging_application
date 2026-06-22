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
	writeTimeout = 5 * time.Second
	readTimeout  = writeTimeout
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
		send:   make(chan []byte),
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
	BroadcastToUser(userID int, event []byte)
	BroadcastError(userID int, ref string, err error)
	BroadcastToConversation(ctx context.Context, conversationID int, event []byte) error
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

func (h *hub) BroadcastToUser(userID int, event []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	conns := h.clients[userID]

	for client := range conns {
		select {
		case client.send <- event:
		default:
			// buffer is full
		}
	}
}

func (h *hub) BroadcastToConversation(ctx context.Context, conversationID int, event []byte) error {
	members, err := h.conversationRepo.GetMembers(ctx, conversationID)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, member := range members {
		conns := h.clients[member.UserID]

		for client := range conns {
			select {
			case client.send <- event:
			default:
				// buffer is full
			}
		}
	}

	return nil
}

func (h *hub) HandleTypingStart(ctx context.Context, client *client, payload ws.TypingPayloadInbound) {
	payloadResp := ws.TypingPayloadOutbound{
		ConversationID: payload.ConversationID,
		UserID:         client.userID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventUserTypingStart, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventUserTypingStart,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventUserTypingStart, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payload.ConversationID, respBytes)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventUserTypingStart, err)
		return
	}
}

func (h *hub) HandleTypingStop(ctx context.Context, client *client, payload ws.TypingPayloadInbound) {
	payloadResp := ws.TypingPayloadOutbound{
		ConversationID: payload.ConversationID,
		UserID:         client.userID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventUserTypingStop, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventUserTypingStop,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventUserTypingStop, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payload.ConversationID, respBytes)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventUserTypingStop, err)
		return
	}
}

func (h *hub) HandleMessageRead(ctx context.Context, client *client, payload ws.MessageReadPayload) {
	err := h.messageService.MarkAsRead(ctx, client.userID, payload.ConversationID, payload.MessageID)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventMessageRead, err)
		return
	}

	payloadResp := ws.MessageSeenPayload{
		UserID:         client.userID,
		MessageID:      payload.MessageID,
		ConversationID: payload.ConversationID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventMessageRead, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventMessageSeen,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventMessageRead, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payloadResp.ConversationID, respBytes)

	if err != nil {
		h.BroadcastError(client.userID, ws.EventMessageRead, err)
		return
	}
}

func (h *hub) BroadcastError(userID int, ref string, err error) {
	payload := ws.ErrorPayload{Ref: &ref}

	if serviceError, ok := errors.AsType[*Error](err); ok {
		payload.Message = serviceError.Message
	} else {
		payload.Message = "something went wrong"
	}

	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return
	}

	resp := ws.Event{
		Type:    "error",
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		return
	}

	h.BroadcastToUser(userID, respBytes)
}

func (h *hub) broadcastOnlineStatus(ctx context.Context, userID int, online bool) {
	conversations, err := h.conversationRepo.GetByUserID(ctx, userID)

	if err != nil {
		return
	}

	var respBytes json.RawMessage

	if online {
		payload := ws.UserOnlinePayload{
			UserID: userID,
		}

		payloadBytes, err := json.Marshal(payload)

		if err != nil {
			return
		}

		resp := ws.Event{
			Type:    ws.EventUserOnline,
			Payload: payloadBytes,
		}

		respBytes, err = json.Marshal(resp)

		if err != nil {
			return
		}
	} else {
		user, err := h.userRepo.UpdateLastSeenAt(ctx, userID)

		if err != nil {
			return
		}

		payload := ws.UserOfflinePayload{
			UserID:     userID,
			LastSeenAt: *user.LastSeenAt,
		}

		payloadBytes, err := json.Marshal(payload)

		if err != nil {
			return
		}

		resp := ws.Event{
			Type:    ws.EventUserOffline,
			Payload: payloadBytes,
		}

		respBytes, err = json.Marshal(resp)

		if err != nil {
			return
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
			h.BroadcastToUser(memberID, respBytes)
		}
	}
}
