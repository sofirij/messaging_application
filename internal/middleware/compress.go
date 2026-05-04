package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
)

func Compress() fiber.Handler {
	return compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	})
}
