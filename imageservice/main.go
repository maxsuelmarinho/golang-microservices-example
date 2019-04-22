package main

import (
	"flag"
	"sync"
	"time"

	"github.com/maxsuelmarinho/golang-microservices-example/common/messaging"

	"github.com/sirupsen/logrus"

	"github.com/maxsuelmarinho/golang-microservices-example/common/config"
	"github.com/maxsuelmarinho/golang-microservices-example/common/logging"
	"github.com/maxsuelmarinho/golang-microservices-example/imageservice/service"
	"github.com/spf13/viper"
)

var appName = "imageservice"

func init() {
	profile := flag.String("profile", "test", "Environment profile")
	configServerUrl := flag.String("config_server_url", "http://config-server:8888", "Address to config server")
	configBranch := flag.String("config_branch", "master", "git branch to fetch configuration from")

	flag.Parse()

	logging.InitializeLogrus(*profile)

	viper.Set("profile", *profile)
	viper.Set("config_server_url", *configServerUrl)
	viper.Set("config_branch", *configBranch)
}

func main() {
	logrus.Infof("Starting %v", appName)

	start := time.Now().UTC()
	config.LoadConfigurationFromBranch(viper.GetString("config_server_url"), appName, viper.GetString("profile"), viper.GetString("config_branch"))
	initializeMessaging()
	go service.StartWebServer(viper.GetString("server_port"))

	logrus.Infof("Started %v in %v", appName, time.Now().UTC().Sub(start))

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func initializeMessaging() {
	if !viper.IsSet("amqp_server_url") {
		panic("No 'amqp_server_url' set in configuration, cannot start")
	}

	service.MessagingClient = &messaging.MessagingClient{}
	service.MessagingClient.ConnectToBroker(viper.GetString("amqp_server_url"))
	service.MessagingClient.Subscribe(viper.GetString("config_event_bus"), "topic", appName, config.HandleRefreshEvent)
}
