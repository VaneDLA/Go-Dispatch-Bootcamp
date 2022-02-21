package repository

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/BernardoGR/Go-Dispatch-Bootcamp/model"
)

func ReadAllCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
			fmt.Errorf("Unable to read input file: %v", err)
			return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	records, err := csvReader.ReadAll()

	if err != nil {
			fmt.Errorf("Unable to parse file: %v", err)
			return nil, err
	}

	return records, nil
}

func AddLine(filePath string, line []string) (error) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
			fmt.Errorf("Unable to read csv file: %v", err)
			return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	err = csvWriter.Write(line)

	if err != nil {
			fmt.Errorf("Unable to write to csv file: %v", err)
			return err
	}

	return nil
}

// Worker Pool implementation
func processLine(wId int, inLines <-chan []string, patients chan<- model.Patient) {
	for line := range inLines {
		log.Printf("Worker ID: %d, line: %s\n", wId, line)
		id, _ := strconv.Atoi(line[0])
		age, _ := strconv.Atoi(line[2])
		patients <- model.Patient {
			ID: id,
		  Name: line[1],
		  Age: age,
		}
	}
}

func ReadCsvFileConcurrent(filePath string) (model.Patients, error) {
	inLines := make(chan []string, 100)
	patients := make(chan model.Patient, 100)

	// turn on 3 go rutines
	for w := 1; w <= 4; w++ {
		go processLine(w, inLines, patients)
	}

	f, err := os.Open(filePath)
	if err != nil {
			fmt.Errorf("Unable to read input file: %v", err)
			return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	if err != nil {
		return nil, err
	}

	i := 1;
	for {
		log.Printf("read new line: %d\n", i)
		line, err := csvReader.Read()
		if i == 1 {
			i += 1
			continue
		}
		if err == io.EOF {
			log.Printf("\n\nEOF found!!\n\n")
			break
		}
		if err != nil {
			fmt.Errorf("Unable to read csv file: %v", err)
			return nil, err
		}
		// send line to worker to be processe
		inLines <- line
		i += 1
  }
	log.Printf("finished reading file")
	close(inLines)

	var result []model.Patient

	for {
		p := <- patients
		log.Printf("patient got from channel: %+v", p)
		result = append(result, p)
		if len(patients) == 0 {
			break
		}
	}
	close(patients)

	log.Printf("total found: %d", len(result))

	log.Printf("WE ARE DONE")

	return result, nil
}
