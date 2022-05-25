package main

import (
	"capstoneProyect/routes"

	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", routes.Index)
	app.Get("/pokemons", routes.GetAll)
	app.Get("/pokemons/:id", routes.GetById)
}

func main() {
	app := fiber.New()

	setupRoutes(app)

	app.Listen("localhost:3000")
}
