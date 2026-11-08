package router

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(r fiber.Router, u *handler.UserHandler, cfg *config.Config) {
	auth := r.Group("/auth")
	auth.Post("/register", u.Register)
	auth.Post("/login", u.Login)
	auth.Post("/logout", middleware.JWT(cfg), u.Logout)
	auth.Post("/refresh", u.RefreshToken)
}
