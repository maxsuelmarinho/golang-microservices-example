package config

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"encoding/json"
	"strings"
	"github.com/spf13/viper"
)

type consumer struct {
	conn *amqp.Connection
	channel *amqp.Channel
	tag string
	done chan error
}

func (c *consumer) Shutdown() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	return <-c.done
}

type updateToken struct {
	Type string `json:"type"`
	Timestamp int `json:"timestamp"`
	OriginService string `json:"originService"`
	DestinationService string `json:"destinationService"`
	ID string `json:"id"`
}

func StartListener(appName string, amqpServer string, exchangeName string) {
	err := newConsumer(amqpServer, exchangeName, "topic", "config-event-queue", exchangeName, appName)
	if err != nil {
		log.Fatalf("%s", err)
	}

	log.Printf("running forever")
	select{}
}

func newConsumer(amqpURI, exchange, exchangeType, queue, key, ctag string) error {
	c := &consumer{
		conn: nil,
		channel: nil,
		tag: ctag,
		done: make(chan error),
	}

	var err error

	log.Printf("dialing %s\n", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	log.Println("got connection, getting channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Printf("go channel, declaring exchange (%s)\n", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,
		exchangeType,
		true, // durable
		false, // delete when complete
		false, // internal
		false, // nowait
		nil, // arguments
	); err != nil {
		return fmt.Errorf("exchange declare: %s", err)
	}

	log.Printf("declared exchange, declaring queue (%s)\n", queue)
	state, err := c.channel.QueueDeclare(
		queue,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // nowait
		nil, // arguments
	)
	if err != nil {
		return fmt.Errorf("Queue declare: %s", err)
	}

	log.Printf("declared queue (%d messages, %d consumers), binding to exchange (key '%s'", state.Messages, state.Consumers, key)

	if err = c.channel.QueueBind(
		queue, // name of the queue
		key, // bindingKey
		exchange, // sourceExchange
		false, // nowait
		nil, // arguments
	); err != nil {
		return fmt.Errorf("Queue bind: %s", err)
	}

	log.Printf("queue bound to exchange, starting consume (consumer tag '%s')", c.tag)
	deliveries, err := c.channel.Consume(
		queue,
		c.tag, //consumerTag
		false, // no ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil, // arguments
	)
	if err != nil {
		return fmt.Errorf("queue consume: %s", err)
	}

	go handle(deliveries, c.done)

	return nil
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf("got %dB consumer: [%v] delivery: [%v] routingKey: [%v] %s", len(d.Body), d.ConsumerTag, d.DeliveryTag, d.RoutingKey, d.Body)

		handleRefreshEvent(d.Body, d.ConsumerTag)
		d.Ack(false)
	}

	log.Printf("handle: deliveries channel closed")
	done <- nil
}

func handleRefreshEvent(body []byte, consumerTag string) {
	updateToken := &updateToken{}
	err := json.Unmarshal(body, updateToken)
	if err != nil {
		log.Printf("error parsing update token: %v", err.Error())
		return
	}

	if strings.Contains(updateToken.DestinationService, consumerTag) {
		log.Println("reloading viper config from spring cloud config server")

		LoadConfigurationFromBranch(
			viper.GetString("config_server_url"),
			consumerTag,
			viper.GetString("profile"),
			viper.GetString("config_branch"),
		)
	}

}