package router

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

func RegisterMessageRoutes(r fiber.Router, h *handler.MessageHandler, cfg *config.Config, userService service.UserService) {
	messages := r.Group("/messages", middleware.JWT(userService, cfg))

	messages.Patch("/:id", h.UpdateBody)
	messages.Delete("/:id", h.SoftDelete)
}
