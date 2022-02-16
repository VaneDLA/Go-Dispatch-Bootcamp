package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BernardoGR/Go-Dispatch-Bootcamp/controller"
	"github.com/BernardoGR/Go-Dispatch-Bootcamp/repository"
	"github.com/BernardoGR/Go-Dispatch-Bootcamp/router"
	"github.com/BernardoGR/Go-Dispatch-Bootcamp/service"
	"github.com/BernardoGR/Go-Dispatch-Bootcamp/usecase"

	"github.com/gorilla/handlers"
)

func main() {
	// app config variables
	dataSource := "csv"
	dataPath := "./resources/patients.csv"

	// Initialize data depending on data source. 
	var raw_data = getRawData(dataSource, dataPath)
	var data = service.ParsePatients(raw_data)

	// create instances for the service, usecase, controller and router
	// injecting the corresponding dependencies to each one of them
	patientService := service.New(data)
	patientUsecase := usecase.New(patientService)
	patientController := controller.New(patientUsecase)
	httpRouter := router.Setup(patientController)

	// Info to set up the server
	// don't use magic naming and magic numbers, there are better ways to do so (viper - covered in another workshop)
	host := "localhost"
	port := 8080

	// create http.Server instance
	server := http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           handlers.LoggingHandler(os.Stdout, httpRouter),
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	log.Printf("starting server in address, %s\n", server.Addr)
	// Start the server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("starting server: %v", err)
	}
}

func getRawData(dataSource, dataPath string) [][]string {
	// use switch in case we have more data sources in the future
	switch dataSource {
	case "csv":
		data, err := repository.ReadCsvFile(dataPath)
		if err != nil {
			log.Println("Error reading csv: ", err)
		}
		return data
	}
	return [][]string{}
}