package router

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(r fiber.Router, h *handler.UserHandler, cfg *config.Config, userService service.UserService) {
	users := r.Group("/users", middleware.JWT(userService, cfg))
	users.Get("/me", h.GetByID)
	users.Delete("/me", h.SoftDelete)
	users.Put("me/avatar", h.UpdateAvatarURL)
	users.Put("/me/username", h.UpdateUsername)
	users.Get("/search", h.SearchByUsername)
}
