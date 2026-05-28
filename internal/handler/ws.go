package handler

import (
	"context"

	"app/internal/model/response"
	"app/internal/service"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

type WSHandler struct {
	hub           service.HubService
	ticketService service.TicketService
}

func NewWSHandler(hub service.HubService, ticketService service.TicketService) *WSHandler {
	return &WSHandler{
		hub:           hub,
		ticketService: ticketService,
	}
}

func (w *WSHandler) Connect(c *websocket.Conn) {
	userID := c.Locals("user_id").(int)
	ctx := context.Background()

	client := service.NewClient(userID, c, ctx)
	w.hub.Register(client)
	go client.WritePump(w.hub)
	client.ReadPump(w.hub)
}

func (t *WSHandler) CreateTicket(c fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	resp, err := t.ticketService.Create(c.Context(), userID)

	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Response[*response.TicketResponse]{Data: resp})
}
