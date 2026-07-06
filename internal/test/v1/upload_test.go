package v1

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"app/internal/model/request"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	filesPath       = "../files/"
	validFilename   = "image.png"
	invalidFilename = "app.exe"
	validFilePath   = filesPath + validFilename
	invalidFilePath = filesPath + invalidFilename
)

func createMultipartRequest(t *testing.T, accessCookie *http.Cookie, fieldname, filename string, content []byte) *http.Request {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	part, err := writer.CreateFormFile(fieldname, filename)
	require.NoError(t, err)

	_, err = io.Copy(part, bytes.NewReader(content))
	require.NoError(t, err)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/upload", &buffer)
	req.AddCookie(accessCookie)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func createMultipartRequestMany(t *testing.T, accessCookie *http.Cookie, fieldname string, filenames []string, content [][]byte) *http.Request {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	for i := range filenames {
		part, err := writer.CreateFormFile(fieldname, filenames[i])
		require.NoError(t, err)

		_, err = io.Copy(part, bytes.NewReader(content[i]))
		require.NoError(t, err)
	}

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/upload-many", &buffer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(accessCookie)

	return req
}

func TestUpload_InvalidFile(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)
	defer clearUploadFolder(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	register(t, app, user1)

	_, accessCookie1, _ := login(t, app, user1)

	content, err := os.ReadFile(invalidFilePath)

	require.NoError(t, err)
	req := createMultipartRequest(t, accessCookie1, "file", invalidFilename, content)

	resp, err := app.Test(req, testCfg)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpload_ValidFile(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)
	defer clearUploadFolder(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	register(t, app, user1)

	_, accessCookie1, _ := login(t, app, user1)

	content, err := os.ReadFile(validFilePath)

	require.NoError(t, err)
	req := createMultipartRequest(t, accessCookie1, "file", validFilename, content)

	resp, err := app.Test(req, testCfg)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUpload_ManyFiles(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)
	defer clearUploadFolder(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	register(t, app, user1)

	_, accessCookie1, _ := login(t, app, user1)

	content, err := os.ReadFile(validFilePath)

	require.NoError(t, err)

	amount := 10
	contents := make([][]byte, amount)
	filenames := make([]string, amount)

	for i := range amount {
		contents[i] = content
		filenames[i] = validFilename
	}

	req := createMultipartRequestMany(t, accessCookie1, "file", filenames, contents)

	resp, err := app.Test(req, testCfg)

	require.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestUpload_TooManyFiles(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)
	defer clearUploadFolder(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	register(t, app, user1)

	_, accessCookie1, _ := login(t, app, user1)

	content, err := os.ReadFile(validFilePath)

	require.NoError(t, err)

	amount := 100
	contents := make([][]byte, amount)
	filenames := make([]string, amount)

	for i := range amount {
		contents[i] = content
		filenames[i] = validFilename
	}

	req := createMultipartRequestMany(t, accessCookie1, "file", filenames, contents)

	resp, err := app.Test(req)

	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpload_ManyFilesOneInvalid(t *testing.T) {
	app := setupApp(t)
	defer truncateTables(t)
	defer clearUploadFolder(t)

	user1 := request.UserAuthRequest{
		Username: "testuser1",
		Password: "password123",
	}

	register(t, app, user1)

	_, accessCookie1, _ := login(t, app, user1)

	content, err := os.ReadFile(validFilePath)

	require.NoError(t, err)

	amount := 9
	contents := make([][]byte, amount)
	filenames := make([]string, amount)

	for i := range amount {
		contents[i] = content
		filenames[i] = validFilename
	}

	// invalid file
	content, err = os.ReadFile(invalidFilePath)

	require.NoError(t, err)

	contents = append(contents, content)
	filenames = append(filenames, invalidFilename)

	req := createMultipartRequestMany(t, accessCookie1, "file", filenames, contents)

	resp, err := app.Test(req, testCfg)

	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func clearUploadFolder(t *testing.T) {
	t.Helper()

	uploadDir := "../../../uploads"
	entries, err := os.ReadDir(uploadDir)

	require.NoError(t, err)
	for _, entry := range entries {
		os.RemoveAll(filepath.Join(uploadDir, entry.Name()))
	}
}
