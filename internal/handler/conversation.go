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

type ConversationHandler struct {
	conversationService service.ConversationService
	hub                 service.HubService
}

func NewConversationHandler(conversationService service.ConversationService, hub service.HubService) *ConversationHandler {
	return &ConversationHandler{
		conversationService: conversationService,
		hub:                 hub,
	}
}

func (h *ConversationHandler) Create(c fiber.Ctx) error {
	var req request.ConversationCreateRequest
	err := c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	userID := c.Locals("user_id").(int)

	resp, err := h.conversationService.Create(c.Context(), userID, req)

	if err != nil {
		return HandleServiceError(c, err)
	}

	// server sent event
	go h.hub.BroadcastToConversation(context.WithoutCancel(c.Context()), userID, ws.EventConversationNew, resp.ID, resp)

	return c.Status(fiber.StatusCreated).JSON(response.Response[*response.ConversationResponse]{Data: resp})
}

func (h *ConversationHandler) GetByUserID(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	resp, err := h.conversationService.GetByUserID(c.Context(), userID)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.JSON(response.Response[[]response.ConversationResponse]{Data: resp})
}

func (h *ConversationHandler) GetByID(c fiber.Ctx) error {
	IDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(IDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	userID := c.Locals("user_id").(int)
	resp, err := h.conversationService.GetByID(c.Context(), userID, conversationID)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.JSON(response.Response[*response.ConversationResponse]{Data: resp})
}

func (h *ConversationHandler) UpdateName(c fiber.Ctx) error {
	IDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(IDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	var req request.ConversationRenameRequest
	err = c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	userID := c.Locals("user_id").(int)
	err = h.conversationService.UpdateName(c.Context(), userID, conversationID, req)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ConversationHandler) UpdateAvatarURL(c fiber.Ctx) error {
	IDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(IDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	var req request.ConversationAvatarRequest
	err = c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	userID := c.Locals("user_id").(int)
	err = h.conversationService.UpdateAvatarURL(c.Context(), userID, conversationID, req)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ConversationHandler) AddMembers(c fiber.Ctx) error {
	IDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(IDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	var req request.ConversationAddMemberRequest
	err = c.Bind().Body(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: err.Error()},
		})
	}

	userID := c.Locals("user_id").(int)
	resp, err := h.conversationService.AddMembers(c.Context(), userID, conversationID, req)

	if err != nil {
		return HandleServiceError(c, err)
	}

	// server sent event
	for _, id := range req.UserIDs {
		payload := ws.MemberAddedPayload{
			ConversationID: conversationID,
			UserID:         id,
		}

		go h.hub.BroadcastToConversation(context.WithoutCancel(c.Context()), userID, ws.EventMemberAdded, conversationID, payload)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Response[[]response.UserResponse]{
		Data: resp,
	})
}

func (h *ConversationHandler) RemoveMember(c fiber.Ctx) error {
	convIDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(convIDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	memIDString := c.Params("uid", "")
	memberID, err := strconv.Atoi(memIDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid member id"},
		})
	}

	userID := c.Locals("user_id").(int)
	err = h.conversationService.RemoveMember(c.Context(), userID, conversationID, memberID)

	if err != nil {
		return HandleServiceError(c, err)
	}

	// server sent event
	payload := ws.MemberRemovedPayload{
		ConversationID: conversationID,
		UserID:         memberID,
	}

	go h.hub.BroadcastToConversation(context.WithoutCancel(c.Context()), userID, ws.EventMemberRemoved, conversationID, payload)
	go h.hub.BroadcastToUser(ws.EventMemberRemoved, memberID, payload)

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ConversationHandler) SoftDelete(c fiber.Ctx) error {
	convIDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(convIDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	userID := c.Locals("user_id").(int)
	err = h.conversationService.SoftDelete(c.Context(), userID, conversationID)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ConversationHandler) ClearMessages(c fiber.Ctx) error {
	convIDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(convIDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	userID := c.Locals("user_id").(int)
	err = h.conversationService.ClearMessages(c.Context(), userID, conversationID)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ConversationHandler) GetMembers(c fiber.Ctx) error {
	convIDString := c.Params("id", "")
	conversationID, err := strconv.Atoi(convIDString)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid conversation id"},
		})
	}

	userID := c.Locals("user_id").(int)

	resp, err := h.conversationService.GetMembers(c.Context(), userID, conversationID)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.JSON(response.Response[[]response.UserResponse]{
		Data: resp,
	})
}
