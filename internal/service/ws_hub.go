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
}

type Hub struct {
	clients          map[int]map[*Client]bool
	register         chan *Client
	unregister       chan *Client
	mu               sync.RWMutex
	conversationRepo repository.ConversationRepository
	messageService   messageService
}

type HubService interface {
	Run()
	Register(client *Client)
	Unregister(client *Client)
	IsOnline(userID int) bool
	BroadcastToUser(userID int, event []byte) error
	BroadcastToConversation(ctx context.Context, conversationID int, event []byte) error
	HandleMessageSend(ctx context.Context, client *Client, payload ws.MessageSendPayload)
	HandleMessageDelete(ctx context.Context, client *Client, payload ws.MessageDeletePayload)
	HandleTypingStart(ctx context.Context, client *Client, payload ws.TypingPayloadInbound)
	HandleTypingStop(ctx context.Context, client *Client, payload ws.TypingPayloadInbound)
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()

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
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) IsOnline(userID int) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.clients[userID]
	return ok
}

func (h *Hub) BroadcastToUser(userID int, event []byte) {
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

func (h *Hub) BroadcastToConversation(ctx context.Context, conversationID int, event []byte) error {
	conversationMembers, err := h.conversationRepo.GetMembers(ctx, conversationID)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, member := range conversationMembers {
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

func (h *Hub) HandleMessageSend(ctx context.Context, client *Client, payload ws.MessageSendPayload) {
	req := request.MessageCreateRequest{
		ReplyToID:   payload.ReplyToID,
		Body:        payload.Body,
		Attachments: payload.Attachments,
	}

	payloadResp, err := h.messageService.Create(ctx, client.UserID, payload.ConversationID, req)

	if err != nil {
		h.sendError(client, ws.EventMessageSend, err)
		return
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.sendError(client, ws.EventMessageSend, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventMessageNew,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.sendError(client, ws.EventMessageSend, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payloadResp.ConversationID, respBytes)

	if err != nil {
		h.sendError(client, ws.EventMessageSend, err)
		return
	}
}

func (h *Hub) HandleMessageDelete(ctx context.Context, client *Client, payload ws.MessageDeletePayload) {
	message, err := h.messageService.SoftDelete(ctx, client.UserID, payload.MessageID)

	if err != nil {
		h.sendError(client, ws.EventMessageDelete, err)
		return
	}

	payloadResp := ws.MessageDeletedPayload{
		MessageID: payload.MessageID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.sendError(client, ws.EventMessageDelete, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventMessageDeleted,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.sendError(client, ws.EventMessageDelete, err)
		return
	}

	err = h.BroadcastToConversation(ctx, message.ConversationID, respBytes)

	if err != nil {
		h.sendError(client, ws.EventMessageDelete, err)
		return
	}
}

func (h *Hub) HandleTypingStart(ctx context.Context, client *Client, payload ws.TypingPayloadInbound) {
	payloadResp := ws.TypingPayloadOutbound{
		ConversationID: payload.ConversationID,
		UserID:         client.UserID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.sendError(client, ws.EventUserTypingStart, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventUserTypingStart,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.sendError(client, ws.EventUserTypingStart, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payload.ConversationID, respBytes)

	if err != nil {
		h.sendError(client, ws.EventUserTypingStart, err)
		return
	}
}

func (h *Hub) HandleTypingStop(ctx context.Context, client *Client, payload ws.TypingPayloadInbound) {
	payloadResp := ws.TypingPayloadOutbound{
		ConversationID: payload.ConversationID,
		UserID:         client.UserID,
	}

	payloadBytes, err := json.Marshal(payloadResp)

	if err != nil {
		h.sendError(client, ws.EventUserTypingStop, err)
		return
	}

	resp := ws.Event{
		Type:    ws.EventUserTypingStop,
		Payload: payloadBytes,
	}

	respBytes, err := json.Marshal(resp)

	if err != nil {
		h.sendError(client, ws.EventUserTypingStop, err)
		return
	}

	err = h.BroadcastToConversation(ctx, payload.ConversationID, respBytes)

	if err != nil {
		h.sendError(client, ws.EventUserTypingStop, err)
		return
	}
}

func (h *Hub) sendError(client *Client, ref string, err error) {
	payload := ws.ErrorPayload{Ref: &ref}

	if serviceError, ok := errors.AsType[*service.Error](err); ok {
		payload.Code = serviceError.Code
		payload.Message = serviceError.Message
	} else {
		payload.Code = "internal_error"
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

	h.BroadcastToUser(client.UserID, respBytes)
}
