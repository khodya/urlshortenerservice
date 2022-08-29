package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getHandler(w, r)
	} else if r.Method == http.MethodPost {
		postHandler(w, r)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "GET")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "POST")
}

func main() {
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(rootHandler),
	}
	server.ListenAndServe()
}
