package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	r *mux.Router
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			fmt.Println("Recovered panic: ", re)

			var err error
			switch re := re.(type) {
			case string:
				err = errors.New(re)
			case error:
				err = re
			default:
				err = errors.New("Unknown error")
			}
			internalErrorHandler(rw, req, err)
		}
	}()

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

	if !checkAuthentication(rw, req) {
		return
	}

	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

type FileServer struct {
	r http.Handler
}

func (s *FileServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if re := recover(); re != nil {
			fmt.Println("Recovered panic: ", re)

			var err error
			switch re := re.(type) {
			case string:
				err = errors.New(re)
			case error:
				err = re
			default:
				err = errors.New("Unknown error")
			}
			internalErrorHandler(rw, req, err)
		}
	}()

	fmt.Printf("%s %s\n", req.Method, req.URL.String())

	// Don't allow external resources
	rw.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; script-src 'self' 'unsafe-inline'")

	// Stop here if its Preflighted OPTIONS request
	if req.Method != "GET" {
		return
	}

	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}
