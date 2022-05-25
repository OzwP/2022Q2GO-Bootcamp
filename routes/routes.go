package routes

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type pokemon struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Power int    `json:"power"`
}

type apiResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous []string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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

func readGeneric(fileName string) map[string]map[string]string {

	f, _ := os.Open(fileName)

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, _ := csvReader.ReadAll()

	elements := make(map[string]map[string]string)
	headers := data[0]

	for _, element := range data[1:] {
		elements[element[0]] = map[string]string{}
		for i, header := range headers {
			elements[element[0]][header] = element[i]
		}
	}

	return elements

}

func makeFile(headers []string, data apiResponse) {

	f, err := os.Create("returnedData.csv")

	if err != nil {
		log.Fatalf("failed creating file: %v", err)
	}

	defer f.Close()

	csvwriter := csv.NewWriter(f)
	csvwriter.Write(headers)
	for _, e := range data.Results {
		csvwriter.Write([]string{e.Name, e.URL})
	}

	csvwriter.Flush()

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
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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

func GetExternal(c *fiber.Ctx) error {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon")
	if err != nil {
		err = fmt.Errorf("request failed %v", err)
		log.Println(err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else if resp.StatusCode != 200 {
		err = fmt.Errorf("could not read from external api, status code: %v", resp.StatusCode)
		log.Println(err)
		return fiber.NewError(resp.StatusCode, err.Error())
	}

	defer resp.Body.Close()

	respB, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error parsing response %v", err)
		log.Println(err)
		return err
	}
	data := apiResponse{}

	if err := json.Unmarshal(respB, &data); err != nil {
		err = fmt.Errorf("error converting response to JSON %v", err)
		log.Println(err)
		return err
	}

	makeFile([]string{"Name", "URL"}, data)

	c.JSON(readGeneric("returnedData.csv"))

	return nil
}
