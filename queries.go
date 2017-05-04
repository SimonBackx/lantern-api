package main

import (
	"encoding/json"
	"fmt"
	//"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
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

	jsonValue, err := json.Marshal(result)
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	fmt.Fprintf(w, "%s", jsonValue)

	/*fmt.Fprintf(w, "[")
	for i, res := range result {
		jsonValue, err := json.Marshal(res)
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
	fmt.Fprintf(w, "]")*/
}

/**
 * /query
 */
func newQueryHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered panic: ", r)
		}
	}()

	if r.Method != "POST" {
		defaultHandler(w, r)
		return
	}

	if !checkAuthentication(w, r) {
		return
	}

	str, err := ioutil.ReadAll(r.Body)

	var query *Query
	err = json.Unmarshal(str, &query)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid query.")
		return
	}

	jsonValue, err := query.MarshalJSONWithoutId()
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	var result interface{}
	err = bson.UnmarshalJSON(jsonValue, &result)

	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	c := mongo.DB("lantern").C("queries")
	if query.Id == nil {
		// new query
		fmt.Println("New query.")
		err = c.Insert(result)

		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}

	} else {
		fmt.Printf("Update query / _id = %v\n", *query.Id)

		err = c.UpdateId(bson.ObjectIdHex(*query.Id), result)
		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}
	}

	/*


		if err != nil {
			internalErrorHandler(w, r, err)
			return
		}*/

	fmt.Fprintf(w, "Success")
}
