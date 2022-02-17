package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/errors"
)

func OpenFile(fileName string) (*os.File, error) {
	csvFile, err := os.Open("data/" + fileName)
	if err != nil {
		log.Printf("Error opening file %v\n", fileName)
		return nil, err
	}
	log.Printf("Successfully opened file data/%v\n", fileName)
	return csvFile, nil
}

func OpenFileForWrite(fileName string) (*os.File, error) {
	csvFile, err := os.OpenFile("data/"+fileName, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		log.Printf("Error opening file %v for writing\n", fileName)
		return nil, err
	}
	log.Printf("Successfully opened file data/%v for writing\n", fileName)
	return csvFile, nil
}

func CloseFile(file *os.File) error {
	if err := file.Close(); err != nil {
		log.Printf("Error closing the file %v\n", file.Name())
		return err
	}
	log.Printf("Successfully closed file %v\n", file.Name())
	return nil
}

func WriteLines(fileName string, data [][]string) error {
	csvFile, err := OpenFileForWrite(fileName)
	if err != nil {
		return &errors.CsvError{
			FileName:    fileName,
			Errors:      []string{err.Error()},
			NumOfErrors: 1,
		}
	}
	defer CloseFile(csvFile)

	w := csv.NewWriter(csvFile)
	w.WriteAll(data)

	if err := w.Error(); err != nil {
		return &errors.CsvError{
			FileName:    fileName,
			Errors:      []string{err.Error()},
			NumOfErrors: 1,
		}
	}

	return nil
}

func ReadLines(fileName string) ([][]string, error) {
	csvFile, err := OpenFile(fileName)
	if err != nil {
		return nil, &errors.CsvError{
			FileName:    fileName,
			Errors:      []string{err.Error()},
			NumOfErrors: 1,
		}
	}
	defer CloseFile(csvFile)

	csvReader := csv.NewReader(csvFile)
	csvError := &errors.CsvError{
		FileName:    fileName,
		Errors:      []string{},
		NumOfErrors: 0,
	}

	result := [][]string{}
	lineCount := 0
	for {
		record, err := csvReader.Read()
		lineCount += 1

		if err == io.EOF {
			break
		}

		if lineCount == 1 {
			continue
		}

		parsedErr, ok := err.(*csv.ParseError)
		if ok && parsedErr.Err == csv.ErrFieldCount {
			error_msg := fmt.Sprintf("Incorrect number of fields. Line: %v. Column: %v.", parsedErr.Line, parsedErr.Column)
			log.Println(error_msg)
			csvError.Errors = append(csvError.Errors, error_msg)
			csvError.NumOfErrors += 1
			result = append(result, record)
			continue
		}

		if err != nil {
			error_msg := fmt.Sprintf("Error reading csv file in line %v: %v\n", lineCount, err.Error())
			log.Println(error_msg)
			csvError.Errors = append(csvError.Errors, error_msg)
			csvError.NumOfErrors += 1
			continue
		}
		result = append(result, record)
	}

	if csvError.NumOfErrors == 0 {
		return result, nil
	}

	return result, csvError
}
