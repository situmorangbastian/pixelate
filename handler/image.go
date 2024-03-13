package handler

import (
	"net/http"
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

	result, err := h.imageService.ConvertPngToJpg(file)
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

	result, err := h.imageService.Resize(file, scale)
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

	result, err := h.imageService.Compress(file)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	}

	return c.SendFile(result)
}
