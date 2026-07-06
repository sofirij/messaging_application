package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CORS(allowedOrigins []string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Accept-Encoding"},
	})
}
