package main

import (
	"capstoneProyect/router"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	router.SetupRoutes(app)

	app.Listen("localhost:3000")
}
