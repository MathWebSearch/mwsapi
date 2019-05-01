package result

import (
	"time"
)

// Result represents an arbirary result
type Result struct {
	Kind string `json:"kind"` // the kind of query this is a result to

	Took           *time.Duration            `json:"took,omitempty"`  // how long the overall query took
	TookComponents map[string]*time.Duration `json:"stats,omitempty"` // how long each of the components took

	Total int64 `json:"total"` // the total number of results
	From  int64 `json:"from"`  // result number this page starts at
	Size  int64 `json:"size"`  // (maximum) number of results in this page

	Variables []*Variable `json:"qvars,omitempty"` // the list of

	HitIDs []int64 `json:"ids,omitempty"`  // the ids of the hits
	Hits   []*Hit  `json:"hits,omitempty"` // the current page of hits
}

// Normalize normalizes a result so that it can be used reproducibly in tests
func (res *Result) Normalize() {
	res.Took = nil
	res.TookComponents = nil
}

// Hit represents a single Hit
type Hit struct {
	ID  string `json:"id,omitempty"`  // the (possibly internal) id of this hit
	URL string `json:"url,omitempty"` // the url of the document returned

	Element *ElasticElement `json:"source,omitempty"` // the raw ElasticSearch element (if any)

	Metadata interface{} `json:"metadata,omitempty"` // arbitrary document meta-data

	Score float64 `json:"score,omitempty"` // score of this hit

	Snippets []string `json:"snippets,omitempty"` // extracts of this hit (if any)
	XHTML    string   `json:"xhtml,omitempty"`    // xhtml source of this hit (if available)

	Math []*MathInfo `json:"math_ids"` // math found within this hit
}
