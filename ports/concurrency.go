package ports

import (
	"capstoneProyect/domain"
	"capstoneProyect/utils"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func WorkerRead(c *fiber.Ctx) error {

	itemType := strings.ToLower(c.Query("type"))
	if itemType != "odd" && itemType != "even" {
		err := fmt.Errorf("not valid type")
		log.Println(err)
		return fiber.NewError(fiber.StatusNotFound, domain.NotFoundMessage)
	}

	items, err1 := strconv.Atoi(c.Query("items", "0"))
	if err1 != nil {
		err1 = fmt.Errorf("number value expected %v", err1)
		log.Println(err1)
		return fiber.NewError(fiber.StatusNotFound, domain.NotFoundMessage)
	}

	items_per_workers, err2 := strconv.Atoi(c.Query("items_per_workers", "0"))
	if err2 != nil {
		err2 = fmt.Errorf("number value expected %v", err2)
		log.Println(err2)
		return fiber.NewError(fiber.StatusNotFound, domain.NotFoundMessage)
	}

	const fName = "/Users/oswaldopacheco/go/wizelineAcademy/2022Q2GO-Bootcamp/data.csv"
	data := utils.ReadFile(fName)

	totalJobs := len(data)

	jobs := make(chan []string)
	results := make(chan []string, items)

	var wg sync.WaitGroup

	nWorkers := int(math.Ceil(float64(items*2) / float64(items_per_workers)))

	wg.Add(nWorkers)

	if items > totalJobs/2 {
		items = totalJobs / 2
	}

	for w := 1; w <= nWorkers; w++ {
		go worker(w, &wg, itemType, items_per_workers, jobs, results)
	}

	go func() {
		for j := items * 2; j > 0; j-- {
			jobs <- data[j]
		}
		close(jobs)
	}()

	wg.Wait()
	close(results)
	responseData := map[string][]string{}

	for resultData := range results {
		responseData[resultData[0]] = resultData
	}

	c.JSON(responseData)
	return nil
}

func worker(id int, wGroup *sync.WaitGroup, iType string, maxJobsPerWorker int, jobs <-chan []string, results chan<- []string) {
	defer wGroup.Done()

	finishedJobs := 0

	for {
		if finishedJobs > maxJobsPerWorker {
			break
		}

		row, ok := <-jobs
		if !ok {
			break
		}

		finishedJobs += 1
		rId, _ := strconv.Atoi(row[0])

		if !(iType == "odd" && rId%2 != 0) && !(iType == "even" && rId%2 == 0) {
			continue
		}
		results <- row

	}
	fmt.Printf("\nWorker %d finished all jobs \n", id)

}
