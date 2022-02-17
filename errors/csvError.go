package errors

import "fmt"

type CsvError struct {
	FileName    string   `json:"file_name"`
	Errors      []string `json:"errors"`
	NumOfErrors int      `json:"num_of_errors"`
}

func (e *CsvError) Error() string {
	errMsg := fmt.Sprintf("Errors in file %v:\n", e.FileName)
	for _, msg := range e.Errors {
		errMsg = fmt.Sprintf("%v%v\n", errMsg, msg)
	}
	return fmt.Sprintf("%vNumber of errors: %v", errMsg, e.NumOfErrors)
}
