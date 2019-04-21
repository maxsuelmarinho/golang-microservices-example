package service

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func StartWebServer(port string) {
	logrus.Infof("Starting HTTP service at %s\n", port)

	r := NewRouter()
	http.Handle("/", r)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		logrus.Errorf("An error occured starting HTTP listener at port %s\n", port)
		logrus.Errorf("Error: %s\n", err.Error())
	}
}
