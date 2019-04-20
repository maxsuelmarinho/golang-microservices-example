package service

import (
	"fmt"
	"net/http"
)

func StartWebServer(port string) {
	r := NewRouter()
	http.Handle("/", r)

	fmt.Printf("Starting HTTP service at %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)

	if err != nil {
		fmt.Printf("An error occured starting HTTP listener at port %s: %s\n", port, err.Error())
	}
}
