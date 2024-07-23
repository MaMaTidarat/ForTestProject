package routes

import (
	"github.com/MaMaTidarat/ForTestProject/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/products", handlers.GetProducts)
}
