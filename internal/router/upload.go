package router

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterUploadRoutes(r fiber.Router, h *handler.UploadHandler, cfg *config.Config) {
	r.Post("/upload", middleware.JWT(cfg), h.Upload)
	r.Post("/upload-many", middleware.JWT(cfg), h.UploadMany)
}
