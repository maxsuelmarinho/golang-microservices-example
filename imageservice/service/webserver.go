package service

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func StartWebServer(port string) {
	r := NewRouter()
	http.Handle("/", r)
	logrus.Infof("Starting HTTP service at %s", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		logrus.Errorf("An error occured starting HTTP listener at port %s: %v", port, err.Error())
	}
}
