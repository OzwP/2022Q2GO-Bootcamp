package ports

import (
	"capstoneProyect/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetAll(c *fiber.Ctx) error {
	pokemons, err := utils.ReadData()
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	c.JSON(pokemons)
	return nil
}
