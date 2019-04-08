package service

import (
	"log"
	"net/http"
)

func StartWebServer(port string) {
	log.Printf("Starting HTTP service at %s\n", port)

	r := NewRouter()
	http.Handle("/", r)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Printf("An error occured starting HTTP listener at port %s\n", port)
		log.Printf("Error: %s\n", err.Error())
	}
}
