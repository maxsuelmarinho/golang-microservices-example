package service

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"ProcessImage",
		"GET",
		"/file/{filename}",
		ProcessImageFromFile,
	},
	Route{
		"GetAccountImage",
		"GET",
		"/accounts/{accountId}",
		GetAccountImage,
	},
	Route{
		"HealthCheck",
		"GET",
		"/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{\"status\":\"UP\"}"))
		},
	},
}
