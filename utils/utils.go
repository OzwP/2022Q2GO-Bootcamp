package utils

import (
	"capstoneProyect/domain"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func ReadFile(fileName string) [][]string {

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

func ReadGeneric(fileName string) map[string]map[string]string {

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

func MakeFile(data domain.ApiResponse) {

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

func ReadData() (map[int]domain.Pokemon, error) {
	const fName = "/Users/oswaldopacheco/go/wizelineAcademy/2022Q2GO-Bootcamp/data.csv"

	data := ReadFile(fName)

	elements := make(map[int]domain.Pokemon)

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

		elements[id] = domain.Pokemon{
			Id:    id,
			Name:  element[1],
			Power: power,
		}

	}
	return elements, nil
}
