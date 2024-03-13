package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/situmorangbastian/pixelate"
)

type imageHttp struct {
	imageService pixelate.ImageService
}

func InitImageHTTP(f *fiber.App, imageService pixelate.ImageService) {
	handler := &imageHttp{imageService}

	f.Post("/convert", handler.convert)
	f.Post("/resize", handler.resize)
	f.Post("/compress", handler.compress)
}

func (h *imageHttp) convert(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if filepath.Ext(file.Filename) != ".png" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid type file"})
	}

	// Open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening uploaded file")
	}
	defer uploadedFile.Close()

	// Create a new file to save the uploaded file
	tempFile, err := os.CreateTemp("./tmp", "uploaded-file-*"+".png")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating temporary file")
	}
	defer tempFile.Close()

	// Copy the file contents to the temporary file
	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error copying file contents")
	}

	result, err := h.imageService.ConvertPngToJpg(tempFile.Name())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.SendFile(result)
}

func (h *imageHttp) resize(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	scale := c.FormValue("scale")
	if scale == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid scale",
		})
	}

	scalePattern := `^\d+:\d+$`
	matched, err := regexp.MatchString(scalePattern, scale)
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	if !matched {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid scale",
		})
	}

	// Open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening uploaded file")
	}
	defer uploadedFile.Close()

	// Create a new file to save the uploaded file
	tempFile, err := os.CreateTemp("./tmp", "uploaded-file-*"+".png")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating temporary file")
	}
	defer tempFile.Close()

	// Copy the file contents to the temporary file
	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error copying file contents")
	}

	result, err := h.imageService.Resize(tempFile.Name(), scale)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.SendFile(result)
}

func (h *imageHttp) compress(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error opening uploaded file")
	}
	defer uploadedFile.Close()

	// Create a new file to save the uploaded file
	tempFile, err := os.CreateTemp("./tmp", "uploaded-file-*"+".png")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error creating temporary file")
	}
	defer tempFile.Close()

	// Copy the file contents to the temporary file
	_, err = io.Copy(tempFile, uploadedFile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error copying file contents")
	}

	result, err := h.imageService.Compress(tempFile.Name())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.SendFile(result)
}
