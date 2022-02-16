package repository

import (
	"encoding/csv"
	"fmt"
	"os"
)

func ReadCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
			fmt.Errorf("Unable to read input file: %v", err)
			return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = 3

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