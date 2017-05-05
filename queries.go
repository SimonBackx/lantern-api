package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"

	"fmt"
	"io/ioutil"
	"net/http"
)

/**
 * /queries
 */
func queriesHandler(w http.ResponseWriter, r *http.Request) {
	c := mongo.DB("lantern").C("queries")
	var result []interface{}
	err := c.Find(nil).Sort("-createdOn").Limit(100).All(&result)
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	jsonValue, err := json.Marshal(result)
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	fmt.Fprintf(w, "%s", jsonValue)
}

/**
 * /query
 */
func newQueryHandler(w http.ResponseWriter, r *http.Request) {

	str, err := ioutil.ReadAll(r.Body)

	var query Query
	err = json.Unmarshal(str, &query)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid query.")
		return
	}

	// todo: verify query

	// Omzetten naar interface{}, anders krijgen we error
	// omdat we de unmarshal functie van bson niet kunnen overschrijven
	var clean interface{}
	jsonValue, err := json.Marshal(query)
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	err = bson.UnmarshalJSON(jsonValue, &clean)
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	c := mongo.DB("lantern").C("queries")
	if query.Id == "" {
		// new query
		fmt.Println("New query.")
		err = c.Insert(clean)

		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}

	} else {
		fmt.Printf("Update query / _id = %s\n", query.Id)

		err = c.UpdateId(query.Id, clean)
		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}
	}

	fmt.Fprintf(w, "Success")
}
