package main

import (
	"fmt"

	"github.com/maxsuelmarinho/golang-microservices-example/quotes/service"
)

var appName = "quote-service"

func main() {
	fmt.Printf("Starting %s on port 8080\n", appName)
	service.StartWebServer("8080")
}
