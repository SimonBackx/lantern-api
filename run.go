package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
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
		// default
		url = "mongodb://lantern:jdgkl6234fsd1DSF08Fsdf@localhost:27017"
	}

	fmt.Printf("Connecting to MongoDB...\n")
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Connected to MongoDB\n")
		return session
	}

	return nil
}

var mongo *mgo.Session

func run(quit chan bool, finished chan bool) {
	defer func() {
		finished <- true
	}()

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("lantrn.xyz"), //your domain here
		Cache:      autocert.DirCache("certs"),           //folder for storing certificates
	}

	// Register test user
	key, found := os.LookupEnv("API_KEY")
	user, found2 := os.LookupEnv("API_USER")
	if found && found2 {
		users[user] = key
		fmt.Printf("Default user: %v=%v\n", user, key)
	} else {
		users["crawler"] = "wQMXWVm4Yab_SKRISvmbWtbWmuMwud7oVRA0JUYThNAYDN8XS8KG4I0uOAOhRUB43rGtbn4VOhyVds-OIseAHwDOUpex0aESRHXz03jbOdSvLRQN-_qTFYqvcU3paXFAEXMz48a7VlB"
		fmt.Printf("Default user not set.\n")
	}

	r := mux.NewRouter()
	// Authenticated requests
	r.HandleFunc("/api/result/{id}", resultHandler).Methods("GET")
	r.HandleFunc("/api/result/{id}/set-category", setResultCategoryHandler).Methods("POST")

	r.HandleFunc("/api/result", newResultHandler).Methods("POST")
	r.HandleFunc("/api/result/{id}", deleteResultHandler).Methods("DELETE")

	r.HandleFunc("/api/results/{queryId}", resultsHandler).Methods("GET")
	r.HandleFunc("/api/results/{queryId}", deleteResultsHandler).Methods("DELETE")
	r.HandleFunc("/api/queries", queriesHandler).Methods("GET")
	r.HandleFunc("/api/query", newQueryHandler).Methods("POST")
	r.HandleFunc("/api/query/{queryId}", deleteQueryHandler).Methods("DELETE")
	r.HandleFunc("/api/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/stats", newStatsHandler).Methods("POST")
	r.HandleFunc("/api/stats", statsHandler).Methods("GET")

	// Not authenticated
	http.HandleFunc("/api/login", loginHandler)
	http.Handle("/api/", &Server{r})
	http.Handle("/", &FileServer{http.FileServer(http.Dir("/etc/lantern/www"))})

	server := &http.Server{
		Addr: ":443",
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	mongo = connectToMongo()

	// Indexen aanmaken
	if mongo != nil {
		index := mgo.Index{
			Key:        []string{"lastFound"},
			Unique:     false,
			DropDups:   false,
			Background: false, // See notes.
			Sparse:     true,  // Enkel als lastFound bestaat, anders niet terug geven als gesorteerd wordt
		}
		c := mongo.DB("lantern").C("results")
		c.EnsureIndex(index)

		index = mgo.Index{
			Key:        []string{"host", "queryId"},
			Unique:     false,
			DropDups:   false,
			Background: false, // See notes.
			Sparse:     true,  // Enkel als host bestaat, anders niet terug geven als gesorteerd wordt
		}
		c.EnsureIndex(index)

		index = mgo.Index{
			Key:        []string{"host", "url", "queryId"},
			Unique:     false,
			DropDups:   false,
			Background: false, // See notes.
			Sparse:     true,  // Enkel als host bestaat, anders niet terug geven als gesorteerd wordt
		}
		c.EnsureIndex(index)

		c = mongo.DB("lantern").C("queries")
		index = mgo.Index{
			Key:        []string{"createdOn"},
			Unique:     false,
			DropDups:   false,
			Background: false, // See notes.
			Sparse:     true,  // Enkel als host bestaat, anders niet terug geven als gesorteerd wordt
		}
		c.EnsureIndex(index)

		go func() {
			err := server.ListenAndServeTLS("", "")
			if err != nil {
				fmt.Println(err.Error())
			}
		}()
	}

	// Wait for finish signal
	<-quit
	if mongo != nil {
		mongo.Close()
	}

	if server != nil {
		server.Shutdown(nil)
	}
}
