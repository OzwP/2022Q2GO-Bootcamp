package main

import (
	// "strconv"
	"capstoneProyect/routes"

	"github.com/gofiber/fiber"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", routes.Index)
	app.Get("/getId/:id", routes.ReadId)
}

func main() {
	app := fiber.New()

	setupRoutes(app)

	app.Listen(3000)
}
