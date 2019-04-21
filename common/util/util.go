package util

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}

	return "127.0.0.1"
}

func GetIPWithPrefix(prefix string) string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "error"
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), prefix) {
			return ipnet.IP.String()
		}
	}

	return "127.0.0.1"
}

func ResolveIPFromHostFile() (string, error) {
	data, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		logrus.Errorf("Problem reading /etc/hosts: %v", err.Error())
		return "", fmt.Errorf("Problem reading /etc/hosts: %v", err.Error())
	}

	lines := strings.Split(string(data), "\n")
	line := lines[len(lines)-1]

	if len(line) < 2 {
		line = lines[len(lines)-2]
	}

	parts := strings.Split(line, "\t")
	return parts[0], nil
}

func HandleSigterm(handleExit func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		handleExit()
		os.Exit(1)
	}()
}
