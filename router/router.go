package router

import (
	"capstoneProyect/ports"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", ports.Index)
	app.Get("/pokemons", ports.GetAll)
	app.Get("/pokemons/:id", ports.GetById)
	app.Get("/external", ports.GetExternal)
	app.Get("/workers", ports.WorkerRead)
}
