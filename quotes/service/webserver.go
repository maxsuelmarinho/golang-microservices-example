package service

import (
	"fmt"
	"log"
	"net/http"
)

func StartWebServer(port string) {
	r := NewRouter()
	http.Handle("/", r)

	log.Printf("Starting HTTP service at %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		log.Printf("An error occured starting HTTP listener at port %s. Error: %s\n", port, err.Error())
	}
}
