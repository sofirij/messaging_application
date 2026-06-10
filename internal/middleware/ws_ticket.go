package middleware

import (
	"app/internal/model/response"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

func TicketAuth(ticketService service.TicketService) fiber.Handler {
	return func(c fiber.Ctx) error {
		ticket := c.Query("ticket")
		if ticket == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: response.ErrorDetail{Message: "missing ticket"},
			})
		}

		userID, err := ticketService.ValidateAndConsume(ticket)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Error: response.ErrorDetail{Message: "ticket expired"},
			})
		}

		c.Locals("user_id", *userID)
		return c.Next()
	}
}
