package handler_test

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/situmorangbastian/pixelate/handler"
	"github.com/situmorangbastian/pixelate/mocks"
)

type funcCall struct {
	Called bool
	Input  []interface{}
	Output []interface{}
}

// setup function to be executed before running tests
func setup() {
	// create a folder tmp
	err := os.Mkdir("./tmp", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(fmt.Errorf("error create folder tmp: %w", err))
	}
}

// teardown function to be executed after running tests
func teardown() {
	if err := os.RemoveAll("tmp"); err != nil {
		panic(fmt.Errorf("error delete folder tmp: %w", err))
	}
}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func TestImageHandler_Convert(t *testing.T) {

	tests := []struct {
		testName               string
		expectedError          bool
		expectedHttpStatusCode int
		imageService           funcCall
		nameFormFile           string
		testFileName           string
	}{
		{
			testName:     "success",
			nameFormFile: "image",
			imageService: funcCall{
				Called: true,
				Input: []interface{}{
					mock.Anything,
				},
				Output: []interface{}{
					"converted.png", nil,
				},
			},
			testFileName: "test.png",
		},
		{
			testName:               "invalid name form file",
			nameFormFile:           "images",
			expectedError:          true,
			expectedHttpStatusCode: http.StatusBadRequest,
			testFileName:           "test.png",
		},
		{
			testName:               "error from service",
			nameFormFile:           "image",
			expectedError:          true,
			expectedHttpStatusCode: http.StatusInternalServerError,
			imageService: funcCall{
				Called: true,
				Input: []interface{}{
					mock.Anything,
				},
				Output: []interface{}{
					"", errors.New("unexpected error"),
				},
			},
			testFileName: "test.png",
		},
		{
			testName:               "invalid extension file",
			nameFormFile:           "image",
			expectedError:          true,
			expectedHttpStatusCode: http.StatusBadRequest,
			testFileName:           "test.jpg",
		},
	}

	app := fiber.New()
	mockImageService := new(mocks.ImageService)
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if test.imageService.Called {
				mockImageService.On("ConvertPngToJpg", test.imageService.Input...).
					Return(test.imageService.Output...).Once()
			}

			fileContent := "file content"
			file := createFormFile("image", test.testFileName, fileContent)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile(test.nameFormFile, file.Filename)
			part.Write([]byte(fileContent))
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/convert", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			handler.InitImageHTTP(app, mockImageService)
			resp, err := app.Test(req)

			mockImageService.AssertExpectations(t)

			if test.expectedError {
				require.Equal(t, test.expectedHttpStatusCode, resp.StatusCode)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestImageHandler_Resize(t *testing.T) {
	tests := []struct {
		testName               string
		scale                  string
		expectedError          bool
		expectedHttpStatusCode int
		imageService           funcCall
		nameFormFile           string
	}{
		{
			testName: "success",
			scale:    "10:10",
			imageService: funcCall{
				Called: true,
				Input: []interface{}{
					mock.Anything, "10:10",
				},
				Output: []interface{}{
					"resized.png", nil,
				},
			},
			nameFormFile: "image",
		},
		{
			testName:               "invalid scale",
			scale:                  "10::10",
			expectedError:          true,
			expectedHttpStatusCode: http.StatusBadRequest,
			nameFormFile:           "image",
		},
		{
			testName:               "empty scale",
			scale:                  "",
			expectedError:          true,
			expectedHttpStatusCode: http.StatusBadRequest,
			nameFormFile:           "image",
		},
		{
			testName:               "error from service",
			scale:                  "10:10",
			expectedError:          true,
			expectedHttpStatusCode: http.StatusInternalServerError,
			imageService: funcCall{
				Called: true,
				Input: []interface{}{
					mock.Anything, "10:10",
				},
				Output: []interface{}{
					"", errors.New("unexpected error"),
				},
			},
			nameFormFile: "image",
		},
		{
			testName:               "invalid name form file",
			scale:                  "10:10",
			expectedError:          true,
			expectedHttpStatusCode: http.StatusBadRequest,
			nameFormFile:           "images",
		},
	}

	app := fiber.New()
	mockImageService := new(mocks.ImageService)
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if test.imageService.Called {
				mockImageService.On("Resize", test.imageService.Input...).
					Return(test.imageService.Output...).Once()
			}

			fileContent := "file content"
			file := createFormFile("image", "test.png", fileContent)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			writer.WriteField("scale", test.scale)
			part, _ := writer.CreateFormFile(test.nameFormFile, file.Filename)
			part.Write([]byte(fileContent))
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/resize", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			handler.InitImageHTTP(app, mockImageService)
			resp, err := app.Test(req)

			mockImageService.AssertExpectations(t)

			if test.expectedError {
				require.Equal(t, test.expectedHttpStatusCode, resp.StatusCode)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestImageHandler_Compress(t *testing.T) {
	tests := []struct {
		testName               string
		expectedError          bool
		expectedHttpStatusCode int
		imageService           funcCall
		nameFormFile           string
	}{
		{
			testName: "success",
			imageService: funcCall{
				Called: true,
				Input: []interface{}{
					mock.Anything,
				},
				Output: []interface{}{
					"compressed.png", nil,
				},
			},
			nameFormFile: "image",
		},
		{
			testName:               "error from service",
			expectedError:          true,
			expectedHttpStatusCode: http.StatusInternalServerError,
			imageService: funcCall{
				Called: true,
				Input: []interface{}{
					mock.Anything,
				},
				Output: []interface{}{
					"", errors.New("unexpected error"),
				},
			},
			nameFormFile: "image",
		},
		{
			testName:               "invalid name form file",
			expectedError:          true,
			expectedHttpStatusCode: http.StatusBadRequest,
			nameFormFile:           "images",
		},
	}

	app := fiber.New()
	mockImageService := new(mocks.ImageService)
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			if test.imageService.Called {
				mockImageService.On("Compress", test.imageService.Input...).
					Return(test.imageService.Output...).Once()
			}

			fileContent := "file content"
			file := createFormFile("image", "test.png", fileContent)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile(test.nameFormFile, file.Filename)
			part.Write([]byte(fileContent))
			writer.Close()

			req := httptest.NewRequest(http.MethodPost, "/compress", body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			handler.InitImageHTTP(app, mockImageService)
			resp, err := app.Test(req)

			mockImageService.AssertExpectations(t)

			if test.expectedError {
				require.Equal(t, test.expectedHttpStatusCode, resp.StatusCode)
				return
			}

			require.NoError(t, err)
		})
	}
}

func createFormFile(fieldName, fileName, fileContent string) *multipart.FileHeader {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fieldName, fileName)
	part.Write([]byte(fileContent))
	writer.Close()

	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	file := &multipart.FileHeader{
		Filename: fileName,
		// Header:   http.Header{"Content-Type": []string{"image/png"}},
		Size: int64(body.Len()),
	}

	return file
}
