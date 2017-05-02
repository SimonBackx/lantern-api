package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Result struct {
	Id        string
	CreatedOn *time.Time
	Url       string
	QueryId   string
	Body      string
}

/**
 * /results/{queryid}
 */
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuthentication(w, r) {
		return
	}

	//c := session.DB(database).C(collection)
	//err := c.Find(query).One(&result)

	fmt.Fprintf(w, "ResultsRequest.")
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
