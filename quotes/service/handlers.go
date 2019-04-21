package service

import (
	"encoding/json"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/maxsuelmarinho/golang-microservices-example/quotes/model"
	"github.com/sirupsen/logrus"
)

var quotes = [...]string{
	"I like a lot of the design decisions they made in the [Go] language. Basically, I like all of them.",
	"In Go, the code does exactly what it says on the page.",
	"Go doesn't implicitly anything.",
	"If I had to describe Go with one word it'd be 'sensible'.",
	"Go will be the server language of the future.",
}

func GetQuote(w http.ResponseWriter, r *http.Request) {
	hostname, _ := os.Hostname()
	addrs, _ := net.LookupHost(hostname)
	addr := ""
	for _, a := range addrs {
		addr = addr + a
	}

	idx := rand.Intn(len(quotes))

	logrus.Infof("Will pick no# %v of the %v quotes\n", idx, len(quotes))
	quote := quotes[idx]
	quoteObject := model.Quote{runtime.GOARCH, runtime.GOOS, hostname + "/" + addr, quote, "EN"}

	data, _ := json.Marshal(quoteObject)
	log.Infof("return string %v\n", string(data))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	data := []byte("{\"status\":\"UP\"}")
	logrus.Infof("return string %v\n", string(data))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
