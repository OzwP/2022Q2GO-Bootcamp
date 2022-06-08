package routes

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

const notFoundMessage = "Item not found"

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

func makeFile(data apiResponse) {

	const fName = "returnedData.csv"

	f, err := os.Create(fName)

	if err != nil {
		log.Fatalf("failed creating file: %v", err)
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalf("failed closong file: %v", err)
		}
	}()

	csvwriter := csv.NewWriter(f)
	csvwriter.Write([]string{"Name", "Url"})
	for _, e := range data.Results {
		csvwriter.Write([]string{e.Name, e.URL})
	}

	csvwriter.Flush()

}

func readData() (map[int]pokemon, error) {
	const fName = "/Users/oswaldopacheco/go/wizelineAcademy/2022Q2GO-Bootcamp/data.csv"

	data := readFile(fName)

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
	const helloMessage = "Hello World!"
	c.Send([]byte(helloMessage))
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
		return fiber.NewError(fiber.StatusNotFound, notFoundMessage)
	}

	if p, ok := pokemons[searchId]; ok {
		c.JSON(p)
	} else {
		return fiber.NewError(fiber.StatusNotFound, notFoundMessage)
	}
	return nil
}

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
	data := apiResponse{}

	if err := json.Unmarshal(respB, &data); err != nil {
		err = fmt.Errorf("error converting response to JSON %v", err)
		log.Println(err)
		return err
	}

	makeFile(data)

	const fName = "returnedData.csv"

	c.JSON(readGeneric(fName))

	return nil
}

func WorkerRead(c *fiber.Ctx) error {

	itemType := c.Query("type")

	items, err1 := strconv.Atoi(c.Query("items", "0"))
	if err1 != nil {
		err1 = fmt.Errorf("number value expected %v", err1)
		log.Println(err1)
		return fiber.NewError(fiber.StatusNotFound, notFoundMessage)
	}

	items_per_workers, err2 := strconv.Atoi(c.Query("items_per_workers", "0"))
	if err2 != nil {
		err2 = fmt.Errorf("number value expected %v", err2)
		log.Println(err2)
		return fiber.NewError(fiber.StatusNotFound, notFoundMessage)
	}

	const fName = "/Users/oswaldopacheco/go/wizelineAcademy/2022Q2GO-Bootcamp/data.csv"
	data := readFile(fName)

	totalJobs := len(data)

	jobs := make(chan []string, totalJobs)
	results := make(chan []string, totalJobs)

	var wg sync.WaitGroup

	nWorkers := items / items_per_workers
	wg.Add(nWorkers)
	totalProcessed := 0
	for w := 1; w <= nWorkers; w++ {
		go worker(w, &wg, itemType, items_per_workers, items, &totalProcessed, jobs, results)
	}

	for j := 1; j < totalJobs; j++ {
		jobs <- data[j]
	}
	close(jobs)

	wg.Wait()

	responseData := map[string][]string{}

	for i := 0; i < items; i++ {
		resultData := <-results
		responseData[resultData[0]] = resultData
	}

	c.JSON(responseData)
	return nil
}

func worker(id int, wGroup *sync.WaitGroup, iType string, maxJobsPerWorker int, totalWanted int, totalProcessed *int, jobs <-chan []string, results chan<- []string) {

	finishedJobs := 1

	for row := range jobs {

		rId, _ := strconv.Atoi(row[0])
		if finishedJobs <= maxJobsPerWorker {
			fmt.Printf("\nWorker %d is working on job number %d", id, finishedJobs)
			if strings.ToLower(iType) == "odd" && rId%2 != 0 {
				finishedJobs += 1
				results <- row
				*totalProcessed += 1
				fmt.Printf("\nWorker %d finished working on job number %d \n\n", id, finishedJobs-1)
			} else if strings.ToLower(iType) == "even" && rId%2 == 0 {
				finishedJobs += 1
				results <- row
				*totalProcessed += 1
				fmt.Printf("\nWorker %d finished working on job number %d \n\n", id, finishedJobs-1)

			}
		} else {
			break
		}

	}
	fmt.Printf("\nWorker %d finished all jobs \n", id)
	wGroup.Done()
}
