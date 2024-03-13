package service

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gofiber/fiber/v2/log"
	"github.com/situmorangbastian/pixelate"
)

type imageService struct{}

func NewImageService() pixelate.ImageService {
	return &imageService{}
}

func (s *imageService) ConvertPngToJpg(file *multipart.FileHeader) (fileName string, err error) {
	src, err := os.Open(file.Filename)
	if err != nil {
		log.Error(err)
		return
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "input-*.png")
	if err != nil {
		log.Error(err)
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, src)
	if err != nil {
		log.Error(err)
		return
	}

	fileName = "converted.jpg"
	cmd := exec.Command("ffmpeg", "-i", tempFile.Name(), fileName)

	// capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Error(stderr.String())
	}

	return
}

func (s *imageService) Resize(file *multipart.FileHeader, scale string) (fileName string, err error) {
	src, err := os.Open(file.Filename)
	if err != nil {
		log.Error(err)
		return
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "input-*"+filepath.Ext(file.Filename))
	if err != nil {
		log.Error(err)
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, src)
	if err != nil {
		log.Error(err)
		return
	}

	ext := filepath.Ext(file.Filename)

	fileName = "resized" + ext
	cmd := exec.Command("ffmpeg", "-i", tempFile.Name(), "-vf", fmt.Sprintf("scale=%s", scale), fileName)

	// capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Error(stderr.String())
	}

	return
}

func (s *imageService) Compress(file *multipart.FileHeader) (fileName string, err error) {
	src, err := os.Open(file.Filename)
	if err != nil {
		log.Error(err)
		return
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "input-*"+filepath.Ext(file.Filename))
	if err != nil {
		log.Error(err)
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, src)
	if err != nil {
		log.Error(err)
		return
	}

	ext := filepath.Ext(file.Filename)

	fileName = "compressed" + ext
	cmd := exec.Command("ffmpeg", "-i", tempFile.Name(), "-crf", "23", fileName)

	// capture standard error
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		log.Error(stderr.String())
	}
	return
}
