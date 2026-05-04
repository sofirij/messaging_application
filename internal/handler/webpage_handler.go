package handler

import (
	staticAssets "app/static"

	"github.com/gofiber/fiber/v3"
)

type WebpageHandler struct{}

func NewWebpageHandler() *WebpageHandler {
	return &WebpageHandler{}
}

func (h *WebpageHandler) ServeHomepage(c fiber.Ctx) error {
	return c.SendFile("dist/index.html", fiber.SendFile{
		FS:       staticAssets.DistFS,
		Compress: true,
	})
}
