package main

import (
	"flag"

	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/dbclient"
	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/service"
	cb "github.com/maxsuelmarinho/golang-microservices-example/common/circuitbreaker"
	"github.com/maxsuelmarinho/golang-microservices-example/common/config"
	"github.com/maxsuelmarinho/golang-microservices-example/common/logging"
	"github.com/maxsuelmarinho/golang-microservices-example/common/messaging"
	"github.com/maxsuelmarinho/golang-microservices-example/common/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var appName = "accountservice"

func init() {
	profile := flag.String("profile", "test", "Environment profile")
	configServerURL := flag.String("configServerUrl", "http://config-server:8888", "Address to config server")
	configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")
	flag.Parse()

	logging.InitializeLogrus(*profile)
	logrus.Infof("Specified configBranch: %s\n", *configBranch)

	viper.Set("profile", *profile)
	viper.Set("config_server_url", *configServerURL)
	viper.Set("config_branch", *configBranch)
}

func initializeMessaging() {
	if !viper.IsSet("amqp_server_url") {
		panic("No 'amqp_server_url' set in configuration, cannot start")
	}

	service.MessagingClient = &messaging.MessagingClient{}
	service.MessagingClient.ConnectToBroker(viper.GetString("amqp_server_url"))
	service.MessagingClient.Subscribe(viper.GetString("config_event_bus"), "topic", appName, config.HandleRefreshEvent)
}

func main() {
	logrus.Infof("Starting %v\n", appName)

	config.LoadConfigurationFromBranch(
		viper.GetString("config_server_url"),
		appName,
		viper.GetString("profile"),
		viper.GetString("config_branch"))

	initializeBoltClient()
	initializeMessaging()
	cb.ConfigureHystrix([]string{"image-service", "quotes-service"}, service.MessagingClient)

	util.HandleSigterm(func() {
		cb.Deregister(service.MessagingClient)
		service.MessagingClient.Close()
	})
	service.StartWebServer(viper.GetString("server_port"))
}

func initializeBoltClient() {
	service.DBClient = &dbclient.BoltClient{}
	service.DBClient.OpenBoltDb()
	service.DBClient.Seed()
}
