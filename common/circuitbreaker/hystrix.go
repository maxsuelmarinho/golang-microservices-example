package circuitbreaker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/maxsuelmarinho/golang-microservices-example/common/messaging"
	"github.com/maxsuelmarinho/golang-microservices-example/common/util"
	"github.com/spf13/viper"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/eapache/go-resiliency/retrier"
	"github.com/sirupsen/logrus"
)

type DiscoveryToken struct {
	Status  string `json:"state"`
	Address string `json:"address"`
}

func init() {
	//logrus.SetOutput(ioutil.Discard)
}

var Client http.Client

var Retries = 3

func CallUsingCircuitBreaker(breakerName string, url string, method string) ([]byte, error) {
	output := make(chan []byte, 1)
	errors := hystrix.Go(breakerName, func() error {
		req, _ := http.NewRequest(method, url, nil)
		err := callWithRetries(req, output)

		return err
	}, func(err error) error {
		logrus.Errorf("In fallback function for breaker %v, error: %v", breakerName, err.Error())
		circuit, _, _ := hystrix.GetCircuit(breakerName)
		logrus.Errorf("Circuit status is: %v", circuit.IsOpen())
		return err
	})

	select {
	case out := <-output:
		logrus.Debugf("Call in breaker %v successful", breakerName)
		return out, nil
	case err := <-errors:
		logrus.Debugf("Got error on channel in breaker %v. Msg: %v", breakerName, err.Error())
		return nil, err
	}
}

func callWithRetries(req *http.Request, output chan []byte) error {
	r := retrier.New(retrier.ConstantBackoff(Retries, 100*time.Millisecond), nil)
	attempt := 0
	err := r.Run(func() error {
		attempt++
		resp, err := Client.Do(req)
		if err == nil && resp.StatusCode < 299 {
			responseBody, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				output <- responseBody
				return nil
			}
			return err
		} else if err == nil {
			err = fmt.Errorf("Status code was %v", resp.StatusCode)
		}

		logrus.Errorf("Retrier failed, attempt %v", attempt)

		return err
	})

	return err
}

func ConfigureHystrix(commands []string, amqpClient messaging.IMessagingClient) {
	for _, command := range commands {
		hystrix.ConfigureCommand(command, hystrix.CommandConfig{
			Timeout:                resolveProperty(command, "Timeout"),
			MaxConcurrentRequests:  resolveProperty(command, "MaxConcurrentRequests"),
			ErrorPercentThreshold:  resolveProperty(command, "ErrorPercentThreshold"),
			RequestVolumeThreshold: resolveProperty(command, "RequestVolumeThreshold"),
			SleepWindow:            resolveProperty(command, "SleepWindow"),
		})
		logrus.Infof("Circuit %v settings: %v", command, hystrix.GetCircuitSettings()[command])
	}

	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "8181"), hystrixStreamHandler)
	logrus.Infoln("Launched hystrixStreamHandler at 8181")

	publishDiscoveryToken(amqpClient)
}

func Deregister(amqpClient messaging.IMessagingClient) {
	ip, err := util.ResolveIPFromHostsFile()
	if err != nil {
		ip = util.GetIPWithPrefix("10.0.")
	}

	token := DiscoveryToken{
		Status:  "DOWN",
		Address: ip,
	}
	bytes, _ := json.Marshal(token)
	amqpClient.PublishOnQueue(bytes, "discovery")
}

func publishDiscoveryToken(amqpClient messaging.IMessagingClient) {
	ip, err := util.ResolveIPFromHostsFile()
	if err != nil {
		ip = util.GetIPWithPrefix("10.0.")
	}
	token := DiscoveryToken{
		Status:  "UP",
		Address: ip,
	}
	bytes, _ := json.Marshal(token)
	go func() {
		for {
			amqpClient.PublishOnQueue(bytes, "discovery")
			time.Sleep(time.Second * 30)
		}
	}()
}

func resolveProperty(command string, prop string) int {
	propertyKey := fmt.Sprintf("hystrix.command.%s.%s", command, prop)
	if viper.IsSet(propertyKey) {
		return viper.GetInt(propertyKey)
	}

	return getDefaultHystrixConfigPropertyValue(prop)
}

func getDefaultHystrixConfigPropertyValue(prop string) int {
	switch prop {
	case "Timeout":
		return hystrix.DefaultTimeout
	case "MaxConcurrentRequests":
		return hystrix.DefaultMaxConcurrent
	case "RequestVolumeThreshold":
		return hystrix.DefaultVolumeThreshold
	case "SleepWindow":
		return hystrix.DefaultSleepWindow
	case "ErrorPercentThreshold":
		return hystrix.DefaultErrorPercentThreshold
	}

	panic("Got unknwon hystrix property: " + prop)
}
