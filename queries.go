package main

import (
	"encoding/json"
	"fmt"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

/**
 * /queries
 */
func queriesHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuthentication(w, r) {
		return
	}

	c := mongo.DB("lantern").C("queries")
	var result []interface{}
	err := c.Find(nil).Sort("-createdOn").Limit(100).All(&result)
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	fmt.Fprintf(w, "[")
	for i, res := range result {
		jsonValue, err := bson.MarshalJSON(res)
		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}

		var query *Query
		err = json.Unmarshal(jsonValue, &query)
		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}

		str, err := query.JSON()
		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}
		if i > 0 {
			fmt.Fprintf(w, ",")
		}
		fmt.Fprintf(w, "%s", str)

	}
	fmt.Fprintf(w, "]")
}
