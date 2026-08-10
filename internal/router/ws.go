package router

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/service"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func RegisterWSRoutes(r fiber.Router, h *handler.WSHandler, ticketService service.TicketService, cfg *config.Config, userService service.UserService) {
	r.Get("/ws/ticket", middleware.JWT(userService, cfg), h.CreateTicket)
	r.Get("/ws", middleware.Upgrade(), middleware.TicketAuth(ticketService), websocket.New(h.Connect))
}
