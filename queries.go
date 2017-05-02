package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

/**
 * /queries
 */
func queriesHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAuthentication(w, r) {
		return
	}

	//c := session.DB(database).C(collection)
	//err := c.Find(query).One(&result)
	query := NewQuery("Testquery", NewOperatorQuery(NewOperatorQuery(NewRegexpQuery("[0-9]+€"), OrOperator, NewRegexpQuery("[0-9]+$")), AndOperator, NewRegexpQuery("\\w+@ey.com")))
	fmt.Println(query.String())
	str, err := query.JSON()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		var original *Query
		err = json.Unmarshal(str, &original)
		if err == nil {
			fmt.Println(original.String())
		} else {
			fmt.Println(err.Error())
		}
		var out bytes.Buffer
		json.Indent(&out, str, "", "\t")
		out.WriteTo(w)

		//fmt.Fprintf(w, "%s", json)
	}
}
