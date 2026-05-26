package service

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"app/internal/config"
	"app/internal/model/response"

	"github.com/google/uuid"
)

// todo: ensure upload size is within app body limit
const maxUploadSize = 99 * 1024 * 1024 // 99MB

var allowedMIMETypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"video/mp4":       true,
	"video/webm":      true,
	"audio/mpeg":      true,
	"audio/ogg":       true,
	"audio/wav":       true,
	"application/pdf": true,
}

type uploadService struct {
	uploadDir string
}

type UploadService interface {
	Upload(file *multipart.FileHeader) (*response.UploadResponse, error)
	UploadMany(files []*multipart.FileHeader) ([]response.UploadResponse, error)
}

func NewUploadService(cfg config.Config) UploadService {
	return &uploadService{
		uploadDir: cfg.UploadDir,
	}
}

func (u *uploadService) Upload(file *multipart.FileHeader) (*response.UploadResponse, error) {
	// file too large
	if file.Size > maxUploadSize {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "file too large",
		}
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	mimeType := http.DetectContentType(buffer[:n])

	// invalid file type
	if !allowedMIMETypes[mimeType] {
		return nil, &Error{
			Code:    ErrCodeBadRequest,
			Message: "invalid file type",
		}
	}

	if _, err = src.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	url := uuid.New().String() + filepath.Ext(file.Filename)
	dstPath := filepath.Join(u.uploadDir, url)

	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		os.Remove(dstPath)
		return nil, err
	}

	return &response.UploadResponse{
		URL:      url,
		Filename: file.Filename,
		Size:     file.Size,
		Type:     mimeType,
	}, nil
}

func (u *uploadService) UploadMany(files []*multipart.FileHeader) ([]response.UploadResponse, error) {
	resp := make([]response.UploadResponse, 0, len(files))

	for _, file := range files {
		r, err := u.Upload(file)

		if err != nil {
			for _, saved := range resp {
				os.Remove(filepath.Join(u.uploadDir, filepath.Base(saved.URL)))
			}
			return nil, err
		}
		resp = append(resp, *r)
	}
	return resp, nil
}
