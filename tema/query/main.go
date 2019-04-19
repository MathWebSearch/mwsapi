package query

import (
	"time"

	"github.com/MathWebSearch/mwsapi/tema"
	"github.com/MathWebSearch/mwsapi/utils"
)

// Query represents a query sent to temasearch
type Query struct {
	// list of possible formula ids
	MathWebSearchIDs []int64

	// Textual content
	Text string
}

// Result represents the result of a Query
type Result struct {
	Total int64 `json:"total"` // the total number of results
	From  int64 `json:"from"`  // result number this page starts at
	Size  int64 `json:"size"`  // (maximum) number of results in this page

	Took         *time.Duration `json:"took"`          // the amount of time the query took to execute, including network latency to elasticsearch
	TookDocument *time.Duration `json:"took_document"` // the amount of time it took for the document query to execute, including it's latency

	Hits []*Hit `json:"hits"` // the current page of hits
}

// Hit represents a single hit of a query
type Hit struct {
	ElasticID string      `json:"id"`       // the id of the document returned
	Metadata  interface{} `json:"metadata"` // the metadata of the hit

	Score    float64  `json:"score"`    // the score this hit achieved
	Snippets []string `json:"snippets"` // the generted snippets

	Math []*ReplacedMath `json:"maths"` // the replaced math elements within this snippet
}

// ReplacedMath represents a single replaced MathExcept
type ReplacedMath struct {
	Source string `json:"source"` // html source code
	ID     string `json:"id"`     // (internal) id
	XPath  string `json:"xpath"`  // the path to the element
}

// RunQuery runs a complete TemaSearch Query
func RunQuery(connection *tema.Connection, q *Query, from int64, size int64) (result *Result, err error) {
	// measure time for this query
	start := time.Now()
	defer func() {
		took := time.Since(start)
		result.Took = &took
	}()

	result = &Result{
		From: from,
		Size: size,
	}

	// run the document query
	res, err := RunDocumentQuery(connection, q, from, size)
	if err != nil {
		return
	}

	result.TookDocument = res.Took
	result.Total = res.Total

	docs := res.Hits

	// prepare running the highlight query in parallel
	highlights := make([]*HighlightResult, len(docs))
	group := utils.NewAsyncGroup()

	for idx, doc := range docs {
		func(idx int, doc *DocumentHit) {
			group.Add(func(_ func(func())) (err error) {
				highlights[idx], err = doc.RunHighlightQuery(connection, q)
				return err
			})
		}(idx, doc)
	}

	// wait for them all to come back
	err = group.Wait()
	if err != nil {
		return
	}

	// and serialize getting the result set back
	result.Hits = make([]*Hit, len(highlights))
	for idx, highlight := range highlights {
		result.Hits[idx], err = NewResult(highlight)
		if err != nil {
			return nil, err
		}
	}

	return
}

// CountQuery counts results subject to a query
func CountQuery(connection *tema.Connection, q *Query) (int64, error) {
	return CountDocumentQuery(connection, q) // the second phase only maps each query 1:1
}

// NewResult generates a new result from a highlight result
func NewResult(highlight *HighlightResult) (result *Hit, err error) {
	result = &Hit{
		ElasticID: highlight.Document.ID,
		Metadata:  highlight.Document.Element.Metadata,

		Score:    *highlight.Hit.Hit.Score,
		Snippets: highlight.Highlights,

		Math: highlight.Math,
	}
	return
}
