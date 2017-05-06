package main

import (
	"encoding/json"
	"fmt"
	"github.com/SimonBackx/lantern-crawler/queries"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

/**
 * /results/{queryId}
 */
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	queryId, found := vars["queryId"]
	if !found {
		internalErrorHandler(w, r, fmt.Errorf("queryId not set"))
		return
	}

	queryIdBson := bson.ObjectIdHex(queryId)

	c := mongo.DB("lantern").C("results")
	var result []queries.Result
	err := c.Find(bson.M{"queryId": queryIdBson}).Select(bson.M{"_id": 1, "queryId": 1, "lastFound": 1, "createdOn": 1, "occurrences": 1, "url": 1, "snippet": 1, "title": 1, "host": 1}).Sort("-lastFound").Limit(500).All(&result)
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
 * /result/{id}
 */
func resultHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, found := vars["id"]
	if !found {
		internalErrorHandler(w, r, fmt.Errorf("id not set"))
		return
	}

	idBson := bson.ObjectIdHex(id)

	c := mongo.DB("lantern").C("results")
	var result queries.Result
	err := c.FindId(idBson).One(&result)
	if err != nil {
		if err == mgo.ErrNotFound {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid id.")
		} else {
			internalErrorHandler(w, r, err)
		}

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
 * /result
 */
func newResultHandler(w http.ResponseWriter, r *http.Request) {
	str, err := ioutil.ReadAll(r.Body)

	var result queries.Result
	err = json.Unmarshal(str, &result)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid result.")
		return
	}

	// todo: verify

	c := mongo.DB("lantern").C("results")
	if result.Id == "" {
		fmt.Println("New result.")
		err = c.Insert(result)

		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}

	} else {
		fmt.Printf("Update result / _id = %v\n", result.Id)

		err = c.UpdateId(result.Id, result)
		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}
	}

	fmt.Fprintf(w, "Success")
}
