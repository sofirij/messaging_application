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
				Code: service.ErrCodeUnauthorized,
			}
		}

		return cfg.JWTPublicKey, nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func verifyAndRefresh(c fiber.Ctx, userService service.UserService, cfg *config.Config) (int, error) {
	token := c.Cookies("access_token")

	claims, err := parseClaims(token, cfg)

	if err != nil {
		accessToken, refreshToken, err := userService.RefreshToken(c.Context(), c.Cookies("refresh_token"))
		if err != nil {
			return 0, err
		}

		c.Cookie(handler.CreateCookie("access_token", accessToken, cfg.AccessTokenDuration))
		c.Cookie(handler.CreateCookie("refresh_token", refreshToken, cfg.RefreshTokenDuration))

		claims, err = parseClaims(accessToken, cfg)

		if err != nil {
			return 0, err
		}
	}

	sub, ok := claims["sub"].(float64)

	if !ok {
		return 0, &service.Error{
			Code: service.ErrCodeUnauthorized,
			Message: "invalid token",
		}
	}

	return int(sub), nil
}

func JWT(userService service.UserService, cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, err := verifyAndRefresh(c, userService, cfg)

		if err != nil {
			return handler.HandleServiceError(c, err)
		}

		c.Locals("user_id", userID)
		return c.Next()
	}
}

func JWTPage(userService service.UserService, cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		_, err := verifyAndRefresh(c, userService, cfg)

		if err != nil {
			return c.Redirect().To("/login")
		}

		return c.Next()
	}
}
