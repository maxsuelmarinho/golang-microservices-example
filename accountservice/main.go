package main

import (
	"fmt"

	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/dbclient"
	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/service"
)

var appName = "accountservice"

func main() {
	fmt.Printf("Starting %v\n", appName)
	initializeBoltClient()
	service.StartWebServer("8080")
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}
