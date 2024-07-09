package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/import-file", importFileHandler)

	fmt.Println("Server started at :3000")
	log.Fatal(app.Listen(":3000"))
}


