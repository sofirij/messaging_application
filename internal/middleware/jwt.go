package middleware

import (

	"app/internal/model/service"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JWT(secret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("access_token")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
		}

		claims := jwt.MapClaims{}

		_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
			// invalid signing method
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, &service.Error{}
			}
			return []byte(secret), nil
		})		

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
		}

		sub, ok := claims["sub"].(float64)

		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
		}

		c.Locals("userID", int(sub))

		return c.Next()
	}
}