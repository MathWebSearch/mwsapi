package result

import (
	"time"
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

// PopulateSubsitutions populates all substiutions within all hits of this result
func (res *Result) PopulateSubsitutions() error {
	for _, hit := range res.Hits {
		if err := hit.PopulateSubsitutions(res); err != nil {
			return err
		}
	}
	return nil
}
