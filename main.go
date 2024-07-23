package main

import (
	"log"

	"poc-app/database"
	"poc-app/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Initialize MongoDB
	database.ConnectDB()

	// Setup routes
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
