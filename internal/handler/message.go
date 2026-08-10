package handler

import (
	"context"
	"strconv"

	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/model/ws"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

type MessageHandler struct {
	messageService service.MessageService
	hub            service.HubService
}

func NewMessageHandler(messageService service.MessageService, hub service.HubService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		hub:            hub,
	}
}

func (m *MessageHandler) Create(c fiber.Ctx) error {
	convIDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(convIDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	var req request.MessageCreateRequest
	err = c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	userID := c.Locals("user_id").(int)
	resp, err := m.messageService.Create(c.Context(), userID, conversationID, req)

	if err != nil {
		return HandleServiceError(c, err)
	}

	// server sent event
	go m.hub.BroadcastToConversation(context.WithoutCancel(c.Context()), userID, ws.EventMessageNew, resp.ConversationID, resp)

	return c.Status(fiber.StatusCreated).JSON(response.Response[*response.MessageResponse]{Data: resp})
}

func (m *MessageHandler) UpdateBody(c fiber.Ctx) error {
	messageIDString := c.Params("id", "")
	messageID, err := strconv.Atoi(messageIDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid message id"},
		})
	}

	var req request.MessageEditRequest
	err = c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	userID := c.Locals("user_id").(int)
	err = m.messageService.UpdateBody(c.Context(), userID, messageID, req)
	if err != nil {
		return HandleServiceError(c, err)
	}

	message, err := m.messageService.GetByID(c.Context(), userID, messageID)

	if err != nil {
		return HandleServiceError(c, err)
	}

	// server sent event
	payload := ws.MessageEditedPayload{
		ConversationID: message.ConversationID,
		MessageID:      message.ID,
		Body:           *message.Body,
	}

	go m.hub.BroadcastToConversation(context.WithoutCancel(c.Context()), userID, ws.EventMessageEdited, message.ConversationID, payload)

	return c.SendStatus(fiber.StatusNoContent)
}

func (m *MessageHandler) SoftDelete(c fiber.Ctx) error {
	messageIDString := c.Params("id", "")
	messageID, err := strconv.Atoi(messageIDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid message id"},
		})
	}

	userID := c.Locals("user_id").(int)
	message, err := m.messageService.SoftDelete(c.Context(), userID, messageID)
	if err != nil {
		return HandleServiceError(c, err)
	}

	// server sent event
	payload := ws.MessageDeletedPayload{
		ConversationID: message.ConversationID,
		MessageID:      message.ID,
	}

	go m.hub.BroadcastToConversation(context.WithoutCancel(c.Context()), userID, ws.EventMessageDeleted, message.ConversationID, payload)

	return c.SendStatus(fiber.StatusNoContent)
}

func (m *MessageHandler) GetByConversationID(c fiber.Ctx) error {
	var before *int
	beforeCursorString := c.Query("before", "")
	if beforeCursorString != "" {
		cursor, err := strconv.Atoi(beforeCursorString)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
				Error: response.ErrorDetail{Message: "invalid before cursor"},
			})
		}
		before = &cursor
	}

	var at *int
	atCursorString := c.Query("at", "")
	if atCursorString != "" {
		cursor, err := strconv.Atoi(atCursorString)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
				Error: response.ErrorDetail{Message: "invalid after cursor"},
			})
		}
		at = &cursor
	}

	if before != nil && at != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "cannot use both 'before' and 'at' query parameter at the same time"},
		})
	}

	limitString := c.Query("limit", "20")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid limit"},
		})
	}

	convIDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(convIDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	userID := c.Locals("user_id").(int)
	resp, err := m.messageService.GetByConversationID(c.Context(), userID, conversationID, before, at, limit)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.JSON(response.Response[*response.PaginatedMessageResponse]{Data: resp})
}
