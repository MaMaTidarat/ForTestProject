package main

import (
    "log"
    "github.com/MaMaTidarat/poc-app/database"
    "github.com/MaMaTidarat/poc-app/routes"
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
