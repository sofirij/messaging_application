package handler

import (
	"app/internal/model/response"
	"app/internal/service"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

const uploadLimit = 10

type UploadHandler struct {
	uploadService service.UploadService
}

func NewUploadHandler(uploadService service.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

func (u *UploadHandler) Upload(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid file"},
		})
	}
	resp, err := u.uploadService.Upload(file)

	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Response[*response.UploadResponse]{Data: resp})
}

func (u *UploadHandler) UploadMany(c fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "invalid files"},
		})
	}

	files, ok := form.File["file"]

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: "missing files"},
		})
	}

	if len(files) > uploadLimit {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Error: response.ErrorDetail{Message: fmt.Sprintf("upload limit of %d files", uploadLimit)},
		})
	}

	resp, err := u.uploadService.UploadMany(files)
	if err != nil {
		return HandleServiceError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response.Response[[]response.UploadResponse]{Data: resp})
}
