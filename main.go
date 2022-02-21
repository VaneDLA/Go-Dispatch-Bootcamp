package main

import (
	"fmt"
	"net/http"

	"github.com/edgarduran/Go-Dispatch-Bootcamp/controllers"
	"github.com/edgarduran/Go-Dispatch-Bootcamp/router"
)

func main() {
	fmt.Printf("Hello\n")
	c := controllers.New()
	router := router.Setup(c)

	fmt.Printf("Listening on port 3030\n")
	http.ListenAndServe(":3030", router)
}
