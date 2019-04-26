package messaging

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type IMessagingClient interface {
	ConnectToBroker(connectionString string)
	Publish(msg []byte, exchangeName string, exchangeType string) error
	PublishOnQueue(msg []byte, queueName string) error
	Subscribe(exchangeName string, exchangeType string, consumerName string, handlerFunc func(amqp.Delivery)) error
	SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery)) error
	Close()
}

type MessagingClient struct {
	conn *amqp.Connection
}

func (m *MessagingClient) ConnectToBroker(connectionString string) {
	if connectionString == "" {
		panic("Cannot initialize connection to broker, connectionString not set.")
	}

	var err error
	m.conn, err = amqp.Dial(fmt.Sprintf("%s/", connectionString))
	if err != nil {
		panic("Failed to connect to AMQP compatible broker at: " + connectionString)
	}
}

func (m *MessagingClient) Publish(body []byte, exchangeName string, exchangeType string) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialized.")
	}
	ch, err := m.conn.Channel()
	defer ch.Close()
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,  // durable
		false, // delete when complete
		false, // internal
		false, // no wait
		nil,   // arguments
	)

	failOnError(err, "Failed to register an Exchange")

	queue, err := ch.QueueDeclare(
		"",    // queue name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)

	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		queue.Name,
		exchangeName, // binding key
		exchangeName, // source exchange
		false,        // no wait
		nil,          // arguments
	)

	failOnError(err, "Failed to bind to the queue")

	err = ch.Publish(
		exchangeName,
		exchangeName, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			Body: body,
		},
	)

	failOnError(err, "Failed to publish a message")

	logrus.Infof("A message was sent: %v\n", body)
	return err
}

func (m *MessagingClient) PublishOnQueue(body []byte, queueName string) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialized.")
	}
	ch, err := m.conn.Channel()
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)

	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	logrus.Infof("A message was sent to queue: %v: %v\n", queueName, string(body))
	return err
}

func (m *MessagingClient) Subscribe(exchangeName string, exchangeType string, consumerName string, handlerFunc func(amqp.Delivery)) error {
	ch, err := m.conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when unused
		false,        // internal
		false,        // no wait
		nil,          // arguments
	)

	failOnError(err, "Failed to register an Exchange")

	logrus.Infof("declared exchange, declaring queue (%s)", "")
	queue, err := ch.QueueDeclare(
		"",    // queue name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)

	failOnError(err, "Failed to register a Queue")

	logrus.Infof("declared queue (%d messages, %d consumers), binding to exchange (key '%s')", queue.Messages, queue.Consumers, exchangeName)

	err = ch.QueueBind(
		queue.Name,   // queue name
		exchangeName, // binding key
		exchangeName, // source exchange
		false,        // no wait
		nil,          // arguments
	)

	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		queue.Name,
		consumerName,
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)

	failOnError(err, "Failed to register a consumer")

	go consumeLoop(msgs, handlerFunc)
	return nil
}

func (m *MessagingClient) SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery)) error {
	ch, err := m.conn.Channel()
	failOnError(err, "Failed to open a channel")

	logrus.Infof("Declaring queue (%s)", queueName)
	queue, err := ch.QueueDeclare(
		queueName,
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)

	failOnError(err, "Failed to register a queue")

	msgs, err := ch.Consume(
		queue.Name,
		consumerName,
		true,  // auto ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)

	failOnError(err, "Failed to register a consumer")

	go consumeLoop(msgs, handlerFunc)
	return nil
}

func (m *MessagingClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery)) {
	for d := range deliveries {
		handlerFunc(d)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		message := fmt.Sprintf("%s: %s", msg, err)
		logrus.Error(message)
		panic(message)
	}
}
