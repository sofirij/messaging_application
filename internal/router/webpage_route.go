package router

import (
	"app/internal/handler"
	staticAssets "app/static"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func registerWebpageRoutes(r fiber.Router, h *handler.WebpageHandler) {

	r.Use("/", static.New("", static.Config{
		FS:         staticAssets.DistFS,
		IndexNames: []string{},
		Browse:     false,
		Compress:   true,
	}))

	r.Get("/*", h.ServeLoginpage)
}
