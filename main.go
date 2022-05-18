package main

import (
	"encoding/csv"
	"os"

	"github.com/gofiber/fiber"
)

func helloWorld(c *fiber.Ctx) {
	c.Send("Hello World")
}

func readCSV(c *fiber.Ctx) {

	f, _ := os.Open("data.csv")

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, _ := csvReader.ReadAll()

	c.Send(data)

}

func main() {
	app := fiber.New()

	app.Get("/", helloWorld)
	app.Get("/readCSV", readCSV)

	app.Listen(3000)
}
