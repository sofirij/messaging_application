package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App, appHost string) {
	app.Use(CORS(appHost))
	app.Use(Logger())
	app.Use(Compress())
}
