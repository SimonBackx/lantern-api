package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

type Operator string

// OperatorQuerys
const (
	AndOperator Operator = "AND"
	OrOperator           = "OR"
)

func (op Operator) String() string {
	if op == AndOperator {
		return "AND"
	}
	return "OR"
}

/// Kan een zoekoperatie uitvoeren op een stuk data.
/// Geeft true / false aan afhankelijk van de match
type QueryAction interface {
	query(str *string) bool
	String() string
}

func UnmarshalQueryAction(b json.RawMessage, destination *QueryAction) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println("error pos 4")
		return err
	}

	var t string
	err = json.Unmarshal(*m["type"], &t)
	if err != nil {
		fmt.Println("error pos 7")
		return err
	}

	switch t {
	case "regexp":
		var r RegexpQuery
		err := json.Unmarshal(b, &r)
		if err != nil {
			fmt.Println("error pos 5")
			return err
		}
		*destination = &r
		return nil
	case "operator":
		var o OperatorQuery
		err := json.Unmarshal(b, &o)
		if err != nil {
			fmt.Println("error pos 6")
			return err
		}
		*destination = &o
		return nil
	}

	return fmt.Errorf("Invalid type")
}

/// Een QueryAction die bestaat uit 2 actions met een operator zoals AND of OR
/// ertussen. Bij and en eerste is false geeft ze meteen false terug
/// Bij OR en eerste is true, geeft het meteen true terug.
type OperatorQuery struct {
	First    QueryAction
	Operator Operator
	Last     QueryAction
}

func (o *OperatorQuery) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})

	m["first"] = o.First
	m["last"] = o.Last
	m["operator"] = o.Operator
	m["type"] = "operator"
	return json.Marshal(m)
}

func (o *OperatorQuery) query(str *string) bool {
	first := o.First.query(str)
	if first {
		if o.Operator == OrOperator {
			return true
		}
	} else {
		if o.Operator == AndOperator {
			return false
		}
	}
	return o.Last.query(str)
}

func (o *OperatorQuery) UnmarshalJSON(b []byte) error {
	// First, deserialize everything into a map of map
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["operator"], &o.Operator)
	if err != nil {
		return err
	}

	err = UnmarshalQueryAction(*objMap["first"], &o.First)
	if err != nil {
		return err
	}

	err = UnmarshalQueryAction(*objMap["last"], &o.Last)
	if err != nil {
		return err
	}
	return nil
}

func (o *OperatorQuery) String() string {
	return "(" + o.First.String() + " " + o.Operator.String() + " " + o.Last.String() + ")"
}

func NewOperatorQuery(first QueryAction, operator Operator, last QueryAction) *OperatorQuery {
	return &OperatorQuery{First: first, Operator: operator, Last: last}
}

/// All supported basic actions
type RegexpQuery struct {
	Regexp *regexp.Regexp
}

func (o *RegexpQuery) MarshalJSON() ([]byte, error) {
	m := make(map[string]string)
	m["regexp"] = o.Regexp.String()
	m["type"] = "regexp"
	return json.Marshal(m)
}

func (o *RegexpQuery) UnmarshalJSON(b []byte) error {
	// First, deserialize everything into a map of map
	var objMap map[string]string
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	o.Regexp, err = regexp.Compile(objMap["regexp"])
	if err != nil {
		return err
	}
	return nil
}

func NewRegexpQuery(str string) *RegexpQuery {
	reg := regexp.MustCompile(str)
	return &RegexpQuery{Regexp: reg}
}

func (a *RegexpQuery) query(str *string) bool {
	return a.Regexp.MatchString(*str)
}

func (q *RegexpQuery) String() string {
	return q.Regexp.String()
}

type ListQuery struct {
}
type TextQuery struct {
}

type Query struct {
	Name      string      `json:"name"`
	CreatedOn *time.Time  `json:"createdOn"`
	Query     QueryAction `json:"root"`
}

func NewQuery(name string, q QueryAction) *Query {
	now := time.Now()
	return &Query{Name: name, CreatedOn: &now, Query: q}
}

func (q *Query) UnmarshalJSON(b []byte) error {
	// First, deserialize everything into a map of map
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*objMap["name"], &q.Name)
	if err != nil {
		fmt.Println("error pos 1")
		return err
	}

	err = json.Unmarshal(*objMap["createdOn"], &q.CreatedOn)
	if err != nil {
		fmt.Println("error pos 2")
		return err
	}

	err = UnmarshalQueryAction(*objMap["root"], &q.Query)
	if err != nil {
		fmt.Println("error pos 3")
		return err
	}
	return nil
}

func (q *Query) String() string {
	return q.Query.String()
}

func (q *Query) JSON() ([]byte, error) {
	return json.Marshal(q)
}
