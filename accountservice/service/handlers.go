package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/dbclient"
	"github.com/maxsuelmarinho/golang-microservices-example/accountservice/model"
	"github.com/maxsuelmarinho/golang-microservices-example/common/messaging"
	"github.com/maxsuelmarinho/golang-microservices-example/common/util"
)

type healthCheckResponse struct {
	Status string `json:"status"`
}

var DBClient dbclient.IBoltClient
var MessagingClient messaging.IMessagingClient
var isHealthy = true
var client = &http.Client{}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	var accountID = mux.Vars(r)["accountId"]
	account, err := DBClient.QueryAccount(accountID)
	if err != nil {
		fmt.Printf("Some error occured serving %s: %s", accountID, err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	account.ServedBy = util.GetIP()

	notifyVIP(account)

	quote, err := getQuote()
	if err == nil {
		account.Quote = quote
	}

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
		fmt.Println("Invalid request to SetHealthyState, allowed values are true or false")
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
}

func getQuote() (model.Quote, error) {
	req, _ := http.NewRequest("GET", "http://quotes-service:8080/api/quote?strength=4", nil)
	resp, err := client.Do(req)

	if err == nil && resp.StatusCode == http.StatusOK {
		quote := model.Quote{}
		bytes, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(bytes, &quote)
		return quote, nil
	}

	return model.Quote{}, fmt.Errorf("Some error")
}

func notifyVIP(account model.Account) {
	if account.ID == "10000" {
		go func(account model.Account) {
			vipNotification := model.VipNotification{
				AccountID: account.ID,
				ReadAt:    time.Now().UTC().String(),
			}
			data, _ := json.Marshal(vipNotification)

			err := MessagingClient.PublishOnQueue(data, "vipQueue")
			if err != nil {
				fmt.Println(err.Error())
			}
		}(account)
	}
}
