package router

import (
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/service"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func RegisterWSRoutes(r fiber.Router, h *handler.WSHandler, ticketService service.TicketService, jwtSecret string) {
	r.Get("/ws/ticket", middleware.JWT(jwtSecret), h.CreateTicket)
	r.Get("/ws", middleware.Upgrade(), middleware.TicketAuth(ticketService), websocket.New(h.Connect))
}
