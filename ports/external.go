package ports

import (
	"capstoneProyect/domain"
	"capstoneProyect/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/gofiber/fiber/v2"
)

func GetExternal(c *fiber.Ctx) error {

	const baseUrl = "pokeapi.co"
	apiBase := "api/v2"
	pokemonEndpoint := "pokemon"

	endpointUrl := "https://" + path.Join(baseUrl, apiBase, pokemonEndpoint)

	resp, err := http.Get(endpointUrl)
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
	data := domain.ApiResponse{}

	if err := json.Unmarshal(respB, &data); err != nil {
		err = fmt.Errorf("error converting response to JSON %v", err)
		log.Println(err)
		return err
	}

	utils.MakeFile(data)

	const fName = "returnedData.csv"

	c.JSON(utils.ReadGeneric(fName))

	return nil
}
