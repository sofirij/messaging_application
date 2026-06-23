package router

import (
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(r fiber.Router, h *handler.UserHandler, jwtSecret string) {
	users := r.Group("/users", middleware.JWT(jwtSecret))
	users.Get("/me", h.GetByID)
	users.Delete("/me", h.SoftDelete)
	users.Put("me/avatar", h.UpdateAvatarURL)
	users.Put("/me/username", h.UpdateUsername)
	users.Get("/search", h.SearchByUsername)
}
