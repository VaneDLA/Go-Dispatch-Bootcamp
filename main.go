package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/controller"
	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/router"
	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/service"
	"github.com/carlos-garibay/Go-Dispatch-Bootcamp/usecase"
	"github.com/gorilla/handlers"
)

func main() {
	ps := service.New(nil)
	uc := usecase.New(ps)
	pc := controller.New(uc)
	httpRouter := router.Setup(pc)

	host := "localhost"
	port := 8080

	server := http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		Handler:           handlers.LoggingHandler(os.Stdout, httpRouter),
		WriteTimeout:      15 * time.Second,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	log.Printf("starting server in address, %s\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("starting server: %v", err)
	}
}
