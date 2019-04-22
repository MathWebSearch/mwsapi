package mws

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Result represents a MathWebSearch result
type Result struct {
	Size  int64 `json:"size,omitempty"` // total size of the resultset (if requested)
	Total int64 `json:"total"`

	Took int64 `json:"time"` // how long the query took, in ms

	Variables []*QueryVariable `json:"qvars"` // list of query variables in the original query
	Hits      []*Hit           `json:"hits"`  // the actual hits of this element
}

// QueryVariable represents a query variable
type QueryVariable struct {
	Name string `json:"name"` // name of the query variable
}

// Hit Represents a single query hit
type Hit struct {
	Formulae []*FormulaeInfo `json:"math_ids"` // Math Elements within this query
	XHTML    string          `json:"xhtml"`    // XHTML Source code of this hit
}

// FormulaeInfo represents a single math excert within an element
type FormulaeInfo struct {
	URL   string `json:"url"`   // url of this element
	XPath string `json:"xpath"` // path of this element
}

// RunRawQuery runs a query
func (conn *Connection) runRawQuery(query *RawQuery) (res *Result, err error) {
	// ensure the query format is json
	query.OutputFormat = "json"

	// turn the query into xml
	b, err := query.ToXML()
	if err != nil {
		return
	}

	// make a request object
	req, err := http.NewRequest("POST", conn.URL(), bytes.NewBuffer(b))
	if err != nil {
		return
	}

	// set some headers
	req.Header.Set("Content-Type", "application/xml")

	// run the request
	resp, err := conn.client.Do(req)
	if err != nil {
		return
	}

	// and get the response
	return NewResult(resp)
}

// NewResult makes a new result given an http response
func NewResult(response *http.Response) (res *Result, err error) {
	defer response.Body.Close()

	// read the body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	// just unmarshal it
	var result Result
	err = json.Unmarshal(body, &result)
	res = &result
	return
}
