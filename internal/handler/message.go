package handler

import (
	"context"
	"encoding/json"
	"strconv"

	"app/internal/model/request"
	"app/internal/model/response"
	"app/internal/model/ws"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

type MessageHandler struct {
	messageService service.MessageService
	hub service.HubService
}

func NewMessageHandler(messageService service.MessageService, hub service.HubService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		hub: hub,
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
		return handleServiceError(c, err)
	}

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
		return handleServiceError(c, err)
	}

	message, err := m.messageService.GetByID(c.Context(), userID, messageID)

	if err != nil {
		return handleServiceError(c, err)
	}

	// server sent event
	func (){
		payload := ws.MessageEditedPayload{
			ConversationID: message.ConversationID,
			MessageID: message.ID,
			Body: *message.Body,
		}

		payloadBytes, err := json.Marshal(payload)

		if err != nil {
			go m.hub.BroadcastError(userID, ws.EventMessageEdited, err)
			return
		}

		event := ws.Event{
			Type: ws.EventMessageEdited,
			Payload: payloadBytes,
		}

		eventBytes, err := json.Marshal(event)

		if err != nil {
			go m.hub.BroadcastError(userID, ws.EventMessageEdited, err)
			return
		}

		go m.hub.BroadcastToConversation(context.WithoutCancel(c.Context()), message.ConversationID, eventBytes)
	}()
	
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
	_, err = m.messageService.SoftDelete(c.Context(), userID, messageID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (m *MessageHandler) GetByConversationID(c fiber.Ctx) error {
	var before *int

	cursorString := c.Query("before", "")
	if cursorString != "" {
		cursor, err := strconv.Atoi(cursorString)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
				Error: response.ErrorDetail{Message: "invalid before cursor"},
			})
		}
		before = &cursor
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
	resp, err := m.messageService.GetByConversationID(c.Context(), userID, conversationID, before, limit)

	if err != nil {
		return handleServiceError(c, err)
	}

	// paginated message response is already arranged as a response.Response
	return c.JSON(resp)
}
