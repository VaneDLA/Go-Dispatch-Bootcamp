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