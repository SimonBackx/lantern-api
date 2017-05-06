package main

import (
	"encoding/json"
	"fmt"
	"github.com/SimonBackx/lantern-crawler/queries"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

var cachedNewResults = make(map[string]int, 0)

func IncreaseResultCount(id bson.ObjectId) {
	str := id.String()
	count, found := cachedNewResults[str]
	if !found {
		cachedNewResults[str] = 1
	} else {
		cachedNewResults[str] = count + 1
	}
}

func DecreaseResultCount(id bson.ObjectId, c int) {
	str := id.String()
	count, found := cachedNewResults[str]
	if found {
		if c > count || c == -1 {
			delete(cachedNewResults, str)
		} else {
			cachedNewResults[str] = count - c
		}
	}
}

func SetResultCount(id bson.ObjectId, count int) {
	str := id.String()
	cachedNewResults[str] = count
}

/**
 * /queries
 */
func queriesHandler(w http.ResponseWriter, r *http.Request) {
	c := mongo.DB("lantern").C("queries")
	var result []bson.M
	err := c.Find(nil).Sort("-createdOn").Limit(100).All(&result)
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	for _, query := range result {
		query["results"] = 0

		id, found := query["_id"]
		if !found {
			continue
		}
		s, ok := id.(bson.ObjectId)
		if !ok {
			continue
		}

		count, found := cachedNewResults[s.String()]
		if found {
			query["results"] = count
		}
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
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Request interrupted")
		return
	}

	var query queries.Query
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

	fmt.Fprintf(w, "ok")
}

/**
 * DELETE /query/{queryid}
 */
func deleteQueryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	queryId, found := vars["queryId"]
	if !found {
		internalErrorHandler(w, r, fmt.Errorf("queryId not set"))
		return
	}

	queryIdBson := bson.ObjectIdHex(queryId)
	c := mongo.DB("lantern").C("queries")
	err := c.RemoveId(queryIdBson)
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	// results deleten
	resultsCollection := mongo.DB("lantern").C("results")
	resultsCollection.RemoveAll(bson.M{"queryId": queryIdBson})
	DecreaseResultCount(queryIdBson, -1)

	fmt.Fprintf(w, "ok")
}
