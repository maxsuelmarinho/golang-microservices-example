package main

import (
	"flag"
	"fmt"

	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/config"
	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/dbclient"
	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/service"
	"github.com/spf13/viper"
)

var appName = "accountservice"

func init() {
	profile := flag.String("profile", "test", "Environment profile")
	configServerURL := flag.String("configServerUrl", "http://config-server:8888", "Address to config server")
	configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")
	flag.Parse()

	fmt.Printf("Specified configBranch: %s\n", configBranch)

	viper.Set("profile", *profile)
	viper.Set("config_server_url", *configServerURL)
	viper.Set("config_branch", *configBranch)
}

func main() {
	fmt.Printf("Starting %v\n", appName)

	config.LoadConfigurationFromBranch(
		viper.GetString("config_server_url"),
		appName,
		viper.GetString("profile"),
		viper.GetString("config_branch"))

	initializeBoltClient()

	go config.StartListener(appName, viper.GetString("amqp_server_url"), viper.GetString("config_event_bus"))

	service.StartWebServer(viper.GetString("server_port"))
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}
