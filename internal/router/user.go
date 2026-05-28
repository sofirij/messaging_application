package router

import (
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(r fiber.Router, h *handler.UserHandler, jwtSecret string) {
	auth := r.Group("/api/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/logout", h.Logout)
	auth.Post("/refresh", h.RefreshToken)

	users := r.Group("/api/users", middleware.JWT(jwtSecret))
	users.Get("/api/me", h.GetByID)
	users.Delete("/me", h.SoftDelete)
	users.Put("me/avatar", h.UpdateAvatarURL)
	users.Put("/me/username", h.UpdateUsername)
	users.Get("/search", h.SearchByUsername)
}