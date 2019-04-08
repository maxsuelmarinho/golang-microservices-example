package main

import (
	"fmt"

	service "github.com/maxsuelmarinho/golang-microservices-example/accountservice/service"
)

var appName = "accountservice"

func main() {
	fmt.Printf("Starting %v\n", appName)
	service.StartWebServer("8080")
}
