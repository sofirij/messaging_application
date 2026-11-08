package middleware

import (
	"app/internal/config"
	"app/internal/handler"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func parseClaims(token string, cfg *config.Config) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		// invalid signing method
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, &service.Error{
				Message: "unexpected signing method",
				Code:    service.ErrCodeUnauthorized,
			}
		}

		return cfg.JWTPublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func verify(c fiber.Ctx, cfg *config.Config) (int, error) {
	token := c.Cookies("access_token")

	claims, err := parseClaims(token, cfg)

	if err != nil {
		return 0, err
	}

	sub, ok := claims["sub"].(float64)

	if !ok {
		return 0, &service.Error{
			Code:    service.ErrCodeUnauthorized,
			Message: "invalid auth token",
		}
	}

	return int(sub), nil
}

func JWT(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, err := verify(c, cfg)

		if err != nil {
			return handler.HandleServiceError(c, &service.Error{
				Message: "invalid auth token",
				Code: service.ErrCodeUnauthorized,
			})
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}
