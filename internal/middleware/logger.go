package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func Logger() fiber.Handler {
	return logger.New(logger.Config{
		LoggerFunc: func (c fiber.Ctx, data *logger.Data, cfg *logger.Config) error {
			accessToken := c.Cookies("access_token", "nil")
			refreshToken := c.Cookies("refresh_token", "nil")

			_, err := fmt.Fprintf(cfg.Stream, "%s | %s | %s | %s | %s -> %d\n", 
			time.Now().Format(time.TimeOnly), accessToken, refreshToken, c.Method(), c.Path(), c.Response().StatusCode())
			
			return err
		},
	})
}
