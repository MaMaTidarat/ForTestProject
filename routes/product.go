package routes

import (
	"github.com/MaMaTidarat/poc-app/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/products", handlers.GetProducts)
}
