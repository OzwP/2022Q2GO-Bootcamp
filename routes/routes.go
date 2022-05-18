package routes

import (
	"encoding/csv"
	_ "fmt"
	"os"

	"github.com/gofiber/fiber"
)

func Index(c *fiber.Ctx) {
	c.Send("Hello World")
}

func ReadId(c *fiber.Ctx) {

	f, _ := os.Open("data.csv")

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, _ := csvReader.ReadAll()

	elements := make(map[string]map[string]string)
	headers := data[0]

	id := c.Params("id")

	for _, element := range data[1:] {
		elements[element[0]] = map[string]string{}
		for i, header := range headers {
			elements[element[0]][header] = element[i]
		}
	}

	if id != "" {
		c.JSON(elements[id])
	} else {
		c.JSON(elements)
	}

}
