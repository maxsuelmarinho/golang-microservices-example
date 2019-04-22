package service

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/maxsuelmarinho/golang-microservices-example/common/messaging"
	"github.com/sirupsen/logrus"
)

var MessagingClient messaging.IMessagingClient

func ProcessImage(w http.ResponseWriter, r *http.Request) {
	sourceImage, _, err := image.Decode(r.Body)
	if err != nil {
		writeServerError(w, err.Error())
		return
	}
	writeAndReturn(w, sourceImage)
}

func GetAccountImage(w http.ResponseWriter, r *http.Request) {
	data := []byte("http://image-service:7777/file/cake.jpg")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func ProcessImageFromFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var filename = vars["filename"]
	logrus.Infof("Serving image for account: %s", filename)

	fImg, err := os.Open(fmt.Sprintf("testimages/%s", filename))
	defer fImg.Close()
	if err != nil {
		logrus.Errorf("Error: %v", err.Error())
		return
	}

	sourceImage, _, err := image.Decode(fImg)
	if err != nil {
		writeServerError(w, err.Error())
		return
	}

	writeAndReturn(w, sourceImage)
}

func writeAndReturn(w http.ResponseWriter, sourceImage image.Image) {
	buf := new(bytes.Buffer)
	err := Sepia(sourceImage, buf)

	if err != nil {
		logrus.Errorf("Error: %v", err.Error())
		writeServerError(w, err.Error())
		return
	}
	outputData := buf.Bytes()

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(outputData)))
	w.WriteHeader(http.StatusOK)
	w.Write(outputData)
}

func writeServerError(w http.ResponseWriter, msg string) {
	logrus.Error(msg)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", msg)))
}
