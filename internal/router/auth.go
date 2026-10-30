package router

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(r fiber.Router, u *handler.UserHandler, cfg *config.Config, userService service.UserService) {
	auth := r.Group("/auth")
	auth.Post("/register", u.Register)
	auth.Post("/login", u.Login)
	auth.Post("/logout", middleware.JWT(userService, cfg), u.Logout)
	auth.Post("/refresh", middleware.JWT(userService, cfg), u.RefreshToken)
}
