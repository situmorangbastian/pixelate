package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/situmorangbastian/pixelate/handler"
)

func main() {
	fiberApp := fiber.New()

	handler.InitImageHTTP(fiberApp)

	// Start server
	go func() {
		if err := fiberApp.Listen(":1111"); err != nil {
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
}
