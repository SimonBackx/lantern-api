package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"net/http"
	"os"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Invalid request.")
}

func internalErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Internal error.")
	fmt.Println(err.Error())
}

func connectToMongo() *mgo.Session {
	url, found := os.LookupEnv("MONGO_URL")

	if !found {
		fmt.Println("MONGO_URL not set in environment variables.")
	} else {
		fmt.Printf("MONGO_URL = %s\n", url)
		session, err := mgo.Dial(url)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Connected to MongoDB\n")
			return session
		}
	}
	return nil
}

var mongo *mgo.Session

func run(quit chan bool, finished chan bool) {
	defer func() {
		finished <- true
	}()

	// Register test user
	key, found := os.LookupEnv("API_KEY")
	user, found2 := os.LookupEnv("API_USER")
	if found && found2 {
		users[user] = key
		fmt.Printf("Default user: %v=%v\n", user, key)
	} else {
		fmt.Printf("Default user not set.\n")
	}

	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/result/", resultHandler)
	http.HandleFunc("/result", newResultHandler)
	http.HandleFunc("/results/", resultsHandler)
	http.HandleFunc("/queries", queriesHandler)

	go func() {
		server.ListenAndServe()
	}()

	mongo = connectToMongo()

	// Wait for finish signal
	<-quit
	if mongo != nil {
		mongo.Close()
	}

	server.Shutdown(nil)
}
