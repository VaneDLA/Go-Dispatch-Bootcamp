package main

import (
	"Go-Dispatch-Bootcamp/src/csv"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

var csvFilePath = "./test/sample.csv"
var records []csv.DataLine

func main() {
	if len(os.Args) == 2 {
		csvFilePath = os.Args[1]
	}

	records = csv.Parse(csvFilePath)

	router := gin.Default()
	router.GET("/data", all)
	router.GET("/data/:id", byId)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Unable to start API server", err)
	}
}

func all(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, records)
}

func byId(c *gin.Context) {
	id := c.Param("id")
	for _, a := range records {
		if a.Id == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Id not found"})
}
