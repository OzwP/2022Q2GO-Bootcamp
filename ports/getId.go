package ports

import (
	"capstoneProyect/domain"
	"capstoneProyect/utils"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetById(c *fiber.Ctx) error {
	pokemons, err := utils.ReadData()
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	searchId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		err = fmt.Errorf("impossible to convert searchId to int: \" %v \" %v", c.Params("id"), err)
		log.Println(err)
		return fiber.NewError(fiber.StatusNotFound, domain.NotFoundMessage)
	}

	if p, ok := pokemons[searchId]; ok {
		c.JSON(p)
	} else {
		return fiber.NewError(fiber.StatusNotFound, domain.NotFoundMessage)
	}
	return nil
}
