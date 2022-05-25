package routes

import (
	"encoding/csv"
	"fmt"
	_ "fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type pokemon struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Power int    `json:"power"`
}

func readFile(fileName string) [][]string {

	f, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	return data

}

func readData() (map[int]pokemon, error) {
	data := readFile("data.csv")

	elements := make(map[int]pokemon)

	for ix, element := range data[1:] {
		id, err := strconv.Atoi(element[0])
		if err != nil {
			err = fmt.Errorf("impossible to parse id from line: %v %v", ix, err)
			return nil, err
		}

		power, err := strconv.Atoi(element[2])
		if err != nil {
			err = fmt.Errorf("impossible to parse power from line: %v %v", ix, err)
			return nil, err
		}

		elements[id] = pokemon{
			id,
			element[1],
			power,
		}

	}
	return elements, nil
}

func Index(c *fiber.Ctx) error {
	c.Send([]byte("Hello World"))
	return nil
}

func GetAll(c *fiber.Ctx) error {
	pokemons, err := readData()
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	c.JSON(pokemons)
	return nil
}

func GetById(c *fiber.Ctx) error {
	pokemons, err := readData()
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	searchId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		err = fmt.Errorf("impossible to convert searchId to int: \" %v \" %v", c.Params("id"), err)
		log.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if p, ok := pokemons[searchId]; ok {
		c.JSON(p)
	} else {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}
	return nil
}
