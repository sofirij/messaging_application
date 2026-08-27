package handler

import (
	"errors"
	"log"

	"app/internal/model/response"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

func HandleServiceError(c fiber.Ctx, err error) error {
	log.Println(err)
	if serviceErr, ok := errors.AsType[*service.Error](err); ok {
		errorDetail := response.ErrorDetail{
			Message: serviceErr.Message,
		}
		resp := response.ErrorResponse{
			Error: errorDetail,
		}

		switch serviceErr.Code {
		case service.ErrCodeUnauthorized:
			resp.Error.Code = fiber.StatusUnauthorized
			return c.Status(fiber.StatusUnauthorized).JSON(resp)
		case service.ErrCodeForbidden:
			resp.Error.Code = fiber.StatusForbidden
			return c.Status(fiber.StatusForbidden).JSON(resp)
		case service.ErrCodeNotFound:
			resp.Error.Code = fiber.StatusNotFound
			return c.Status(fiber.StatusNotFound).JSON(resp)
		case service.ErrCodeConflict:
			resp.Error.Code = fiber.StatusConflict
			return c.Status(fiber.StatusConflict).JSON(resp)
		default:
			resp.Error.Code = fiber.StatusBadRequest
			return c.Status(fiber.StatusBadRequest).JSON(resp)
		}
	} else {
		errorDetail := response.ErrorDetail{
			Code: fiber.StatusInternalServerError,
			Message: "something went wrong",
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Error: errorDetail,
		})
	}
}
