package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/dbclient"
	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/model"
	cb "github.com/maxsuelmarinho/golang-microservices-example/common/circuitbreaker"
	"github.com/maxsuelmarinho/golang-microservices-example/common/messaging"
	"github.com/maxsuelmarinho/golang-microservices-example/common/util"
	"github.com/sirupsen/logrus"
)

type healthCheckResponse struct {
	Status string `json:"status"`
}

var DBClient dbclient.IBoltClient
var MessagingClient messaging.IMessagingClient
var isHealthy = true
var client = &http.Client{}
var fallbackQuote = model.Quote{
	Language: "en",
	ServedBy: "circuit-breaker",
	Text:     "May the source be with you, always.",
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	var accountID = mux.Vars(r)["accountId"]
	account, err := DBClient.QueryAccount(accountID)
	if err != nil {
		logrus.Errorf("Some error occured serving %s: %s", accountID, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	account.ServedBy = util.GetIP()

	notifyVIP(account)

	quote, err := getQuote()
	if err == nil {
		account.Quote = quote
	}
	account.ImageURL = getImageURL(accountID)

	data, _ := json.Marshal(account)
	writeJsonResponse(w, http.StatusOK, data)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	dbUP := DBClient.Check()
	if dbUP && isHealthy {
		data, _ := json.Marshal(healthCheckResponse{Status: "UP"})
		writeJsonResponse(w, http.StatusOK, data)
	} else {
		data, _ := json.Marshal(healthCheckResponse{Status: "Database unaccessible"})
		writeJsonResponse(w, http.StatusServiceUnavailable, data)
	}
}

func SetHealthyState(w http.ResponseWriter, r *http.Request) {
	var state, err = strconv.ParseBool(mux.Vars(r)["state"])

	if err != nil {
		logrus.Errorf("Invalid request to SetHealthyState, allowed values are true or false")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isHealthy = state
	w.WriteHeader(http.StatusOK)
}

func writeJsonResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(status)
	w.Write(data)
}

func init() {
	var transport http.RoundTripper = &http.Transport{
		DisableKeepAlives: true,
	}
	client.Transport = transport
	cb.Client = *client
}

func getQuote() (model.Quote, error) {
	body, err := cb.CallUsingCircuitBreaker("quotes-service", "http://quotes-service:8080/api/quote?strength=4", "GET")

	if err == nil {
		quote := model.Quote{}
		json.Unmarshal(body, &quote)
		return quote, nil
	}

	return fallbackQuote, nil
}

func notifyVIP(account model.Account) {
	if account.ID == "10000" {
		go func(account model.Account) {
			vipNotification := model.VipNotification{
				AccountID: account.ID,
				ReadAt:    time.Now().UTC().String(),
			}
			data, _ := json.Marshal(vipNotification)

			err := MessagingClient.PublishOnQueue(data, "vip_queue")
			if err != nil {
				logrus.Error(err.Error())
			}
		}(account)
	}
}

func getImageURL(accountID string) string {
	body, err := cb.CallUsingCircuitBreaker("image-service", fmt.Sprintf("http://image-service:8080/accounts/%s", accountID), "GET")
	if err == nil {
		return string(body)
	} else {
		return "http://path.to.placeholder"
	}
}
