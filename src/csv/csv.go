package csv

import (
	"encoding/csv"
	"log"
	"os"
)

type DataLine struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	IPAddress string `json:"ip_address"`
}

func Parse(filePath string) []DataLine {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal("Unable to read close file "+filePath, err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	var csvData []DataLine

	for _, value := range records {
		newLine := DataLine{Id: value[0], FirstName: value[1], LastName: value[2], Email: value[3], Gender: value[4], IPAddress: value[5]}
		csvData = append(csvData, newLine)
	}

	return csvData
}
