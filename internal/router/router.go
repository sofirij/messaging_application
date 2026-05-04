package router

import (
	"app/internal/handler"

	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App) {

	handlers := handler.Setup()

	registerWebpageRoutes(app, handlers.WebpageHandler)
}
