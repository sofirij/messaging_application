package router

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterMessageRoutes(r fiber.Router, h *handler.MessageHandler, cfg *config.Config) {
	messages := r.Group("/messages", middleware.JWT(cfg))

	messages.Patch("/:id", h.UpdateBody)
	messages.Delete("/:id", h.SoftDelete)
}
