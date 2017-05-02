package main

import (
	"fmt"
	"net/http"
)

var users = make(map[string]string, 0)

func checkAuthentication(w http.ResponseWriter, r *http.Request) bool {
	return true

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
