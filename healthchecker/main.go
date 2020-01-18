package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	host := flag.String("host", "http://127.0.0.1:8080", "host to check")

	resp, err := http.Get(fmt.Sprintf("%s/health", *host))
	if err != nil || resp.StatusCode != 200 {
		os.Exit(1)
	}
	os.Exit(0)
}
