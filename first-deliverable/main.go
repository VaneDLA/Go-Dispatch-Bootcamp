package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"os"
)

type IdRecord struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/data/{id}", fetchById)
	r.Get("/data", fetchAll)
	http.ListenAndServe(":3000", r)
}

func fetchAll(w http.ResponseWriter, r *http.Request) {
	records, err := readCsvFile()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Failed with %s", err)))
		return
	}
	render.JSON(w, r, records)
}

func fetchById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	record, err := findCsvFile(id)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf("Failed with %s", err)))
		return
	}
	render.JSON(w, r, record)
}

func readCsvFile() ([]IdRecord, error) {
	csvfile, err := os.Open("data.csv")

	var csvData = []IdRecord{}
	if err != nil {
		log.Fatal("Unable to read input file", err)
		return csvData, err
	}
	defer csvfile.Close()

	csvReader := csv.NewReader(csvfile)
	records, err := csvReader.ReadAll()
	if err != nil {
		return csvData, err
	}
	for _, value := range records {
		newLine := IdRecord{ID: value[0], Name: value[1]}
		csvData = append(csvData, newLine)
	}
	return csvData, nil
}

func findCsvFile(id string) (IdRecord, error) {
	csvfile, err := os.Open("data.csv")
	if err != nil {
		log.Fatal("Unable to read input file", err)
	}
	defer csvfile.Close()

	csvReader := csv.NewReader(csvfile)
	searchedLine := IdRecord{}
	records, err := csvReader.ReadAll()
	if err != nil {
		return searchedLine, err
	}
	for _, value := range records {
		if value[0] == id {
			searchedLine = IdRecord{ID: value[0], Name: value[1]}
			break
		}
	}
	if searchedLine.ID == "" {
		return searchedLine, errors.New("Record not found!")
	}

	return searchedLine, nil
}
