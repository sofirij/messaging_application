package handler

import (
	"app/internal/model/response"
	"app/internal/service"

	"github.com/gofiber/fiber/v3"
)

type uploadHandler struct {
	uploadService service.UploadService
}

func NewUploadHandler(uploadService service.UploadService) *uploadHandler {
	return &uploadHandler{
		uploadService: uploadService,
	}
}

func (u *uploadHandler) Upload(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid file"},
		})
	}
	resp, err := u.uploadService.Upload(file)

	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Response[*response.UploadResponse]{Data: resp})
}

func (u *uploadHandler) UploadMany(c fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid files"},
		})
	}

	files, ok := form.File["files"]

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "missing files"},
		})
	}

	resp, err := u.uploadService.UploadMany(files)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Response[[]response.UploadResponse]{Data: resp})
}
