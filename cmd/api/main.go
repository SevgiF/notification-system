package main

import (
	"github.com/SevgiF/notification-system/internal/core/notification"
	"github.com/SevgiF/notification-system/pkg/database/mysql"
	"github.com/gofiber/fiber/v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Initialize MySQL connection
	dbManager := mysql.NewMySQLConnectionManager()
	db := dbManager.DB
	if err := db.Ping(); err != nil {
		log.Fatalf("MYSQL veri tabanına erişilemiyor : %v", err)
	}

	app := fiber.New()
	notification.SetupNotification(app, db)

	// Setup graceful shutdown
	// Create channel for shutdown signals
	shutdownChan := make(chan os.Signal, 1)
	// Listen for SIGINT and SIGTERM signals
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Listen(":80"); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Block until a signal is received
	<-shutdownChan
	log.Println("Shutting down server...")

	// Shutdown all resources gracefully
	// 1. Shutdown HTTP server
	if err := app.Shutdown(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// 2. Close database connection
	if err := db.Close(); err != nil {
		log.Printf("Database shutdown error: %v", err)
	}

	log.Println("Server gracefully stopped")
}
