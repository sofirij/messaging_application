package router

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

func RegisterUploadRoutes(r fiber.Router, h *handler.UploadHandler, cfg *config.Config, userService service.UserService) {
	r.Post("/upload", middleware.JWT(userService, cfg), h.Upload)
	r.Post("/upload-many", middleware.JWT(userService, cfg), h.UploadMany)
}
