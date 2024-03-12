package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type imageHttp struct {
}

func InitImageHTTP(f *fiber.App) {
	handler := &imageHttp{}

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

	src, err := file.Open()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "input-*.png")
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, src); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	outputFileName := "output.jpg"
	cmd := exec.Command("ffmpeg", "-i", tempFile.Name(), outputFileName)

	// Capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Error(stderr.String())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.SendFile(outputFileName)
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

	src, err := file.Open()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "input-*"+filepath.Ext(file.Filename))
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, src); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	ext := filepath.Ext(file.Filename)

	outputFileName := "resized" + ext
	cmd := exec.Command("ffmpeg", "-i", tempFile.Name(), "-vf", fmt.Sprintf("scale=%s", scale), outputFileName)

	// Capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Error(stderr.String())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.SendFile(outputFileName)
}

func (h *imageHttp) compress(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	src, err := file.Open()
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "input-*"+filepath.Ext(file.Filename))
	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}
	defer os.Remove(tempFile.Name())

	if _, err := io.Copy(tempFile, src); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	ext := filepath.Ext(file.Filename)

	outputFileName := "compressed" + ext
	cmd := exec.Command("ffmpeg", "-i", tempFile.Name(), "-vf", "-crf", "23", outputFileName)

	// Capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Error(stderr.String())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.SendFile(outputFileName)
}
