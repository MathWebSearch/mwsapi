package result

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/MathWebSearch/mwsapi/utils"
)

// Result represents a Query Result which can be from either MathWebSearch, ElasticSearch or a combined TemaSearch Result
// too unmarshal
type Result struct {
	Kind Kind `json:"kind"` // the kind of query this is a result to

	Took           *time.Duration            `json:"took,omitempty"`  // how long the overall query took
	TookComponents map[string]*time.Duration `json:"stats,omitempty"` // how long each of the components took

	Total int64 `json:"total"` // the total number of results
	From  int64 `json:"from"`  // result number this page starts at
	Size  int64 `json:"size"`  // (maximum) number of results in this page

	Variables []*QueryVariable `json:"qvars,omitempty"` // the list of

	HitIDs []int64 `json:"ids,omitempty"`  // the ids of the hits
	Hits   []*Hit  `json:"hits,omitempty"` // the current page of hits
}

// Kind represents the kind of result being returned
type Kind string

const (
	// EmptyKind is a ressult in response to the empty query
	EmptyKind Kind = ""
	// MathWebSearchKind is a result in represonse to a MathWebSearch query
	MathWebSearchKind Kind = "mwsd"
	// ElasticDocumentKind is a result returned from the document query
	ElasticDocumentKind Kind = "elastic-document"
	// ElasticKind is a result returned from the document + highlight query
	ElasticKind Kind = "elastic"
	// TemaSearchKind is a result returned from a combined elastic + MathWebSearch Query
	TemaSearchKind Kind = "tema"
)

// Normalize normalizes a result so that it can be used reproducibly in tests
func (res *Result) Normalize() {
	res.Took = nil
	res.TookComponents = nil
}

// UnmarshalMWS unmarshals given a response from an html server
func (res *Result) UnmarshalMWS(response *http.Response) error {
	defer response.Body.Close() // close the body when done
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// read the adapted data from the json
	var r struct {
		Total    int64 `json:"total"`
		TookInMS int   `json:"time"`

		Variables []*QueryVariable `json:"qvars"`

		MathWebSearchIDs []int64 `json:"ids,omitempty"`
		Hits             []*Hit  `json:"hits,omitempty"`
	}
	if err := json.Unmarshal(responseBytes, &r); err != nil {
		return err
	}

	// store the time taken
	took := time.Duration(r.TookInMS) * time.Millisecond
	res.Took = &took

	res.Kind = MathWebSearchKind
	res.Total = r.Total

	// TODO: Update this
	res.TookComponents = map[string]*time.Duration{
		"mwsd": &took,
	}

	res.Size = utils.MaxInt64(int64(len(r.Hits)), int64(len(r.MathWebSearchIDs)))

	res.Variables = r.Variables
	res.HitIDs = r.MathWebSearchIDs
	res.Hits = r.Hits

	return nil
}
