package main

import (
	"flag"
	"fmt"

	"github.com/maxsuelmarinho/golang-microservices-example/common/config"
	"github.com/maxsuelmarinho/golang-microservices-example/common/messaging"
	"github.com/maxsuelmarinho/golang-microservices-example/common/util"
	"github.com/maxsuelmarinho/golang-microservices-example/vipservice/service"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var appName = "vipservice"
var messagingClient messaging.IMessagingClient

func init() {
	configServerUrl := flag.String("config_server_url", "http://config-server:8888", "Address to config server")
	profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
	configBranch := flag.String("config_branch", "master", "git branch to fetch configuration from")

	flag.Parse()

	viper.Set("profile", *profile)
	viper.Set("config_server_url", *configServerUrl)
	viper.Set("config_branch", *configBranch)
}

func main() {
	fmt.Printf("Starting %s...\n", appName)

	config.LoadConfigurationFromBranch(viper.GetString("config_server_url"), appName, viper.GetString("profile"), viper.GetString("config_branch"))
	initializeMessaging()

	util.HandleSigterm(func() {
		if messagingClient != nil {
			messagingClient.Close()
		}
	})

	service.StartWebServer(viper.GetString("server_port"))
}

func onMessage(delivery amqp.Delivery) {
	fmt.Printf("Got a message: %v\n", string(delivery.Body))
}

func initializeMessaging() {
	if !viper.IsSet("amqp_server_url") {
		panic("No 'amqp_server_url' set in configuration, cannot start")
	}
	messagingClient = &messaging.Messagingclient{}
	messagingClient.ConnectToBroker(viper.GetString("amqp_server_url"))

	err := messagingClient.SubscribeToQueue("vip_queue", appName, onMessage)
	failOnError(err, "Could not start subscribe to vip_queue")

	err = messagingClient.Subscribe(viper.GetString("config_event_bus"), "topic", appName, config.HandleRefreshEvent)
	failOnError(err, fmt.Sprintf("Could not start subscribe to \"%s\" topic\n", viper.GetString("config_event_bus")))
}

func failOnError(err error, msg string) {
	if err != nil {
		message := fmt.Sprintf("%s: %s\n", msg, err)
		fmt.Println(message)
		panic(message)
	}
}
