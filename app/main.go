package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/situmorangbastian/pixelate/handler"
	"github.com/situmorangbastian/pixelate/service"
	"github.com/spf13/viper"
)

func main() {
	// create a folder tmp
	err := os.Mkdir("tmp", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(fmt.Errorf("error create folder tmp: %w", err))
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error config file: %w", err))
	}

	servicePort := viper.GetInt("service.port")
	if servicePort <= 0 {
		panic("invalid service port")
	}

	imageService := service.NewImageService()

	fiberApp := fiber.New()

	handler.InitImageHTTP(fiberApp, imageService)

	// Start server
	go func() {
		if err := fiberApp.Listen(fmt.Sprintf(":%s", strconv.Itoa(servicePort))); err != nil {
			log.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err := fiberApp.ShutdownWithTimeout(10 * time.Second); err != nil {
		log.Fatal(err)
	}

	// remove output file
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	extToRemoves := []string{".png", ".jpg"}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		for _, extToRemove := range extToRemoves {
			if !info.IsDir() && filepath.Ext(path) == extToRemove {
				err := os.Remove(path)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll("tmp"); err != nil {
		panic(fmt.Errorf("error delete folder tmp: %w", err))
	}
}
