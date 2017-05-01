package main

import (
	"fmt"
	"net/http"
	"strconv"
)

var users = make(map[string]string, 0)

func checkAuthentication(w http.ResponseWriter, r *http.Request) bool {
	user := r.Header.Get("X-API-USER")
	key := r.Header.Get("X-API-KEY")

	if user == "" || key == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Not authorized.")
		return false
	}

	foundKey, found := users[user]
	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Not authorized.")
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	if foundKey != key {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Not authorized.")
		return false
	}

	return true
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Invalid request.")
}

/**
 * /result/{id}
 */
func resultHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuthentication(w, r) {
		return
	}

	_id := r.URL.Path[len("/result/"):]
	id, err := strconv.Atoi(_id)
	if err != nil {
		// invalid number
		defaultHandler(w, r)
		return
	}

	fmt.Fprintf(w, "ResultRequest. id=%v", id)
}

/**
 * /result
 */
func newResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		defaultHandler(w, r)
		return
	}

	if !checkAuthentication(w, r) {
		return
	}

	fmt.Fprintf(w, "New result")
}

func run(quit chan bool, finished chan bool) {
	defer func() {
		finished <- true
	}()

	// Register test user
	users["default"] = "default"

	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/result/", resultHandler)
	http.HandleFunc("/result", newResultHandler)

	go func() {
		server.ListenAndServe()
	}()

	// Wait for finish signal
	<-quit

	server.Shutdown(nil)
}
