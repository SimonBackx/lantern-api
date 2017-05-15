package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"io/ioutil"
	"net/http"
)

var users = make(map[string]string, 0)

func checkAuthentication(w http.ResponseWriter, r *http.Request) bool {
	user := r.Header.Get("X-API-USER")
	key := r.Header.Get("X-API-KEY")

	if user == "" || key == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Not authorized.")
		return false
	}

	foundKey, found := users[user]
	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Not authorized.")
		return false
	}

	if foundKey != key {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Not authorized.")
		return false
	}

	return true
}

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SavedLoginCredentials struct {
	Id       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string        `json:"username" bson:"username"`
	Password string        `json:"password" bson:"password"`
}

/**
 * POST /register (authenticated!)
 */
func registerHandler(w http.ResponseWriter, r *http.Request) {
	str, err := ioutil.ReadAll(r.Body)

	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	var credentials LoginCredentials
	err = json.Unmarshal(str, &credentials)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid credentials.")
		return
	}

	c := mongo.DB("lantern").C("users")
	var foundUser SavedLoginCredentials
	err = c.Find(bson.M{"username": credentials.Username}).One(&foundUser)

	if err == nil || err != mgo.ErrNotFound {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Username already in use.")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	err = c.Insert(SavedLoginCredentials{Username: credentials.Username, Password: string(hashedPassword)})

	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	fmt.Fprintf(w, "ok")
}

/**
 * POST /login (unauthenticated!)
 */
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	str, err := ioutil.ReadAll(r.Body)

	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	var credentials LoginCredentials
	err = json.Unmarshal(str, &credentials)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid credentials.")
		return
	}

	c := mongo.DB("lantern").C("users")
	var foundUser SavedLoginCredentials
	err = c.Find(bson.M{"username": credentials.Username}).One(&foundUser)

	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid credentials.")
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(credentials.Password))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid credentials.")

		// todo: add rate limit for this user
		return
	}

	// Tokens uitschrijven
	key, err := GenerateRandomString(512)

	if err != nil {
		internalErrorHandler(w, r, err)
		return
	}

	user := foundUser.Id.Hex()
	users[user] = key
	fmt.Fprintf(w, "{\"user\": \"%s\", \"key\": \"%s\"}", user, key)

}
