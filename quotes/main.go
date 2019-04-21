package main

import (
	"flag"

	"github.com/maxsuelmarinho/golang-microservices-example/common/logging"
	"github.com/maxsuelmarinho/golang-microservices-example/quotes/service"
	"github.com/sirupsen/logrus"
)

var appName = "quote-service"

func init() {
	profile := flag.String("profile", "test", "Environment profile")
	flag.Parse()
	logging.InitializeLogrus(*profile)
}

func main() {
	logrus.Infof("Starting %s on port 8080\n", appName)
	service.StartWebServer("8080")
}
