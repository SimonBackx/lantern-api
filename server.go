package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	r *mux.Router
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s %s\n", req.Method, req.URL.String())

	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-API-USER, X-API-KEY")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}
