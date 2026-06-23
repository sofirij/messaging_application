package router

import (
	"app/internal/handler"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(r fiber.Router, u *handler.UserHandler, jwtSecret string) {
	auth := r.Group("/auth")
	auth.Post("/register", u.Register)
	auth.Post("/login", u.Login)
	auth.Post("/logout", u.Logout)
	auth.Post("/refresh", u.RefreshToken)
}