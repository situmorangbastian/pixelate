package main

import (
	"fmt"
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

	dir, err := os.Getwd() // Get the current working directory
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Target extension to remove
	extToRemoves := []string{".jpg", ".png", ".jpeg"} // Change this to your desired extension

	// Walk through the directory and remove files with the specified extension
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, extToRemove := range extToRemoves {
			if !info.IsDir() && filepath.Ext(path) == extToRemove {
				err := os.Remove(path)
				if err != nil {
					return err
				}
				fmt.Printf("Deleted: %s\n", path)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	}
}
