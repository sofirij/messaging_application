package router

import (
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterMessageRoutes(r fiber.Router, h *handler.MessageHandler, jwtSecret string) {
	messages := r.Group("/messages", middleware.JWT(jwtSecret))

	messages.Patch("/:id", h.UpdateBody)
	messages.Delete("/:id", h.SoftDelete)
}
