package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"

	"github.com/maxsuelmarinho/golang-microservices-example/geltfail/aggregator"
	"github.com/maxsuelmarinho/golang-microservices-example/geltfail/transformer"
	"github.com/sirupsen/logrus"
)

var authToken = ""
var port *string

func init() {
	data, err := ioutil.ReadFile("token.text")
	if err != nil {
		msg := "Cannot find token.txt that should contain the Loggly token"
		logrus.Errorln(msg)
		panic(msg)
	}
	authToken = string(data)
	port = flag.String("port", "12201", "UDP port for the gelftail")
	flag.Parse()
}

func main() {
	logrus.Println("Starting Gelt-tail server...")

	serverConn := startUDPServer(*port)
	defer serverConn.Close()

	var bulkQueue = make(chan []byte, 1) // buffered channel to put log statements ready for LaaS upload into

	go aggregator.Start(bulkQueue, authToken)
	go listenForLogStatements(serverConn, bulkQueue)

	logrus.Infoln("Started Gelf-tail server")
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func startUDPServer(port string) *net.UDPConn {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
	checkError(err)

	serverConn, err := net.ListenUDP("udp", serverAddr)
	checkError(err)

	return serverAddr
}

func checkError(err error) {
	if err != nil {
		logrus.Errorf("Error: %v", err)
		os.Exit(0)
	}
}

func listenForLogStatements(serverConn *net.UDPConn, bulkQueue chan []byte) {
	buf := make([]byte, 8192) // buffer to store UDP payload into. 8kb
	var item map[string]interface{}

	for {
		n, _, err := serverConn.ReadFromUDP(buf) // Blocks until data becomes available, which is put into the buffer
		if err != nil {
			logrus.Errorf("Problem reading UDP message into buffer: %s", err.Error())
			continue
		}

		err = json.Unmarshal(buf[0:n], &item)
		if err != nil {
			logrus.Errorf("Problem unmarshalling log message into JSON: %s", err.Error())
			continue
		}

		processedLogMessage, err := transformer.ProcessLogStatement(item)

		if err != nil {
			logrus.Errorf("Problem parsing message: %v", string(buf[0:n]))
		} else {
			bulkQueue <- processedLogMessage
		}
		item = nil
	}
}
