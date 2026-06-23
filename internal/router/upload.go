package router

import (
	"app/internal/handler"
	"app/internal/middleware"

	"github.com/gofiber/fiber/v3"
)

func RegisterUploadRoutes(r fiber.Router, h *handler.UploadHandler, jwtSecret string) {
	r.Post("/upload", middleware.JWT(jwtSecret), h.Upload)
	r.Post("/upload-many", middleware.JWT(jwtSecret), h.UploadMany)
}
