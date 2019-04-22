package config

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func (c *consumer) Shutdown() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer logrus.Infof("AMQP shutdown OK")

	return <-c.done
}

type updateToken struct {
	Type               string `json:"type"`
	Timestamp          int    `json:"timestamp"`
	OriginService      string `json:"originService"`
	DestinationService string `json:"destinationService"`
	ID                 string `json:"id"`
}

func StartListener(appName string, amqpServer string, exchangeName string) {
	err := newConsumer(amqpServer, exchangeName, "topic", "config-event-queue", exchangeName, appName)
	if err != nil {
		logrus.Fatalf("%s", err)
	}

	logrus.Infof("running forever")
	select {}
}

func newConsumer(amqpURI, exchange, exchangeType, queue, key, ctag string) error {
	c := &consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	logrus.Infof("dialing %s\n", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	logrus.Infof("got connection, getting channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	logrus.Infof("go channel, declaring exchange (%s)\n", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		true,  // durable
		false, // delete when complete
		false, // internal
		false, // nowait
		nil,   // arguments
	); err != nil {
		return fmt.Errorf("exchange declare: %s", err)
	}

	logrus.Infof("declared exchange, declaring queue (%s)\n", queue)
	state, err := c.channel.QueueDeclare(
		queue,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // nowait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue declare: %s", err)
	}

	logrus.Infof("declared queue (%d messages, %d consumers), binding to exchange (key '%s'", state.Messages, state.Consumers, key)

	if err = c.channel.QueueBind(
		queue,    // name of the queue
		key,      // bindingKey
		exchange, // sourceExchange
		false,    // nowait
		nil,      // arguments
	); err != nil {
		return fmt.Errorf("Queue bind: %s", err)
	}

	logrus.Infof("queue bound to exchange, starting consume (consumer tag '%s')", c.tag)
	deliveries, err := c.channel.Consume(
		queue,
		c.tag, //consumerTag
		false, // no ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("queue consume: %s", err)
	}

	go handle(deliveries, c.done)

	return nil
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		logrus.Infof("got %dB consumer: [%v] delivery: [%v] routingKey: [%v] %s", len(d.Body), d.ConsumerTag, d.DeliveryTag, d.RoutingKey, d.Body)

		HandleRefreshEvent(d)
		d.Ack(false)
	}

	logrus.Infof("handle: deliveries channel closed")
	done <- nil
}

func HandleRefreshEvent(d amqp.Delivery) {
	body := d.Body
	consumerTag := d.ConsumerTag
	updateToken := &updateToken{}
	err := json.Unmarshal(body, updateToken)
	if err != nil {
		logrus.Errorf("error parsing update token: %v", err.Error())
		return
	}

	if strings.Contains(updateToken.DestinationService, consumerTag) {
		logrus.Info("reloading viper config from spring cloud config server")

		LoadConfigurationFromBranch(
			viper.GetString("config_server_url"),
			consumerTag,
			viper.GetString("profile"),
			viper.GetString("config_branch"),
		)
	}

}
