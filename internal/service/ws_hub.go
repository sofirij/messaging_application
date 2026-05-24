package service

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"app/internal/model/request"
	"app/internal/model/service"
	"app/internal/model/ws"
	"app/internal/repository"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	UserID int
	conn   *websocket.Conn
	send   chan []byte
	ctx    context.Context
}

type hub struct {
	clients          map[int]map[*Client]bool
	register         chan *Client
	unregister       chan *Client
	mu               sync.RWMutex
	messageService   MessageService
	conversationRepo repository.ConversationRepository
	userRepo         repository.UserRepository
}

type HubService interface {
	Run()
	Register(client *Client)
	Unregister(client *Client)
	IsOnline(userID int) bool
	BroadcastToUser(userID int, event []byte)
	BroadcastToConversation(ctx context.Context, conversationID int, event []byte) error
	HandleMessageSend(ctx context.Context, client *Client, payload ws.MessageSendPayload)
	HandleMessageDelete(ctx context.Context, client *Client, payload ws.MessageDeletePayload)
	HandleTypingStart(ctx context.Context, client *Client, payload ws.TypingPayloadInbound)
	HandleTypingStop(ctx context.Context, client *Client, payload ws.TypingPayloadInbound)
	HandleMessageRead(ctx context.Context, client *Client, payload ws.MessageReadPayload)
}

func NewHubService(conversationRepo repository.ConversationRepository, userRepo repository.UserRepository, messageService MessageService) HubService {
	return &hub{
		clients:          make(map[int]map[*Client]bool),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		mu:               sync.RWMutex{},
		messageService:   messageService,
		conversationRepo: conversationRepo,
		userRepo:         userRepo,
	}
}

func (h *hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			h.broadcastUserStatus(client.ctx, client.UserID, true)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				delete(clients, client)
				close(client.send)
				if len(clients) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
			h.broadcastUserStatus(client.ctx, client.UserID, false)
		}
	}
}

func (h *hub) Register(client *Client) {
	h.register <- client
}

func (h *hub) Unregister(client *Client) {
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

func (h *hub) HandleMessageSend(ctx context.Context, client *Client, payload ws.MessageSendPayload) {
	req := request.MessageCreateRequest{
		ReplyToID:   payload.ReplyToID,
		Body:        payload.Body,
		Attachments: payload.Attachments,
	}

	payloadResp, err := h.messageService.Create(ctx, client.UserID, payload.ConversationID, req)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageSend, err)
		return
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageSend, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventMessageNew,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageSend, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payloadResp.ConversationID, respBytes)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageSend, err)
		return
	}
}

func (h *hub) HandleMessageDelete(ctx context.Context, client *Client, payload ws.MessageDeletePayload) {
	message, err := h.messageService.SoftDelete(ctx, client.UserID, payload.MessageID)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageDelete, err)
		return
	}

	payloadResp := ws.MessageDeletedPayload{
		MessageID:      payload.MessageID,
		ConversationID: message.ConversationID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageDelete, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventMessageDeleted,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageDelete, err)
		return
	}

	err = h.BroadcastToConversation(ctx, message.ConversationID, respBytes)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageDelete, err)
		return
	}
}

func (h *hub) HandleTypingStart(ctx context.Context, client *Client, payload ws.TypingPayloadInbound) {
	payloadResp := ws.TypingPayloadOutbound{
		ConversationID: payload.ConversationID,
		UserID:         client.UserID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventUserTypingStart, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventUserTypingStart,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventUserTypingStart, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payload.ConversationID, respBytes)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventUserTypingStart, err)
		return
	}
}

func (h *hub) HandleTypingStop(ctx context.Context, client *Client, payload ws.TypingPayloadInbound) {
	payloadResp := ws.TypingPayloadOutbound{
		ConversationID: payload.ConversationID,
		UserID:         client.UserID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventUserTypingStop, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventUserTypingStop,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventUserTypingStop, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payload.ConversationID, respBytes)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventUserTypingStop, err)
		return
	}
}

func (h *hub) HandleMessageRead(ctx context.Context, client *Client, payload ws.MessageReadPayload) {
	err := h.messageService.MarkAsRead(ctx, client.UserID, payload.ConversationID, payload.MessageID)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageRead, err)
		return
	}

	payloadResp := ws.MessageSeenPayload{
		UserID:         client.UserID,
		MessageID:      payload.MessageID,
		ConversationID: payload.ConversationID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageRead, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventMessageSeen,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageRead, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payloadResp.ConversationID, respBytes)

	if err != nil {
		h.broadcastError(client.UserID, ws.EventMessageRead, err)
		return
	}
}

func (h *hub) broadcastError(userID int, ref string, err error) {
	payload := ws.ErrorPayload{Ref: &ref}

	if serviceError, ok := errors.AsType[*service.Error](err); ok {
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

func (h *hub) broadcastUserStatus(ctx context.Context, userID int, online bool) {
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

	// verify user status hasn't changed before broadcasting
	if (!online && h.IsOnline(userID)) || (online && !h.IsOnline(userID)) {
		return
	}

	for memberID := range memberMap {
		if memberID != userID {
			h.BroadcastToUser(memberID, respBytes)
		}
	}
}
