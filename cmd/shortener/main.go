package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	scheme      = "http"
	host        = "localhost"
	port        = ":8080"
	hostAndPort = host + port
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getHandler(w, r)
	} else if r.Method == http.MethodPost {
		postHandler(w, r)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if strings.Count(r.URL.Path, "/") > 1 {
		http.Error(w, "Bad shortened url in path.", http.StatusBadRequest)
		return
	}
	var urlToDecode string
	if strings.HasPrefix(r.URL.Path, "/") {
		urlToDecode = r.URL.Path[1:]
	} else {
		urlToDecode = r.URL.Path
	}
	decodedBytes, err := decode(urlToDecode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	decoded := string(decodedBytes)
	if _, err := url.Parse(decoded); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", decoded)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Request body is empty.", http.StatusBadRequest)
		return
	}
	if _, err := url.Parse(string(body)); err != nil {
		http.Error(w, "Could not parse url from request.", http.StatusBadRequest)
		return
	}
	encodedURL := encode(body)
	w.WriteHeader(http.StatusCreated)
	url := &url.URL{Scheme: scheme, Host: hostAndPort, Path: encodedURL}
	w.Write([]byte(url.String()))
}

func encode(v []byte) string {
	return base64.StdEncoding.EncodeToString(v)
}

func decode(v string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(v)
}

func main() {
	server := &http.Server{
		Addr:    port,
		Handler: http.HandlerFunc(rootHandler),
	}
	server.ListenAndServe()
}
