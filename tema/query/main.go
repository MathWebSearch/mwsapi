package query

import "github.com/MathWebSearch/mwsapi/tema"

// Query represents a query sent to temasearch
type Query struct {
	// list of possible formula ids
	MathWebSearchIDs []int64

	// Textual content
	Text string
}

// Result represents the result of a temasearch query
type Result struct {
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
func RunQuery(connection *tema.Connection, q *Query, from int64, size int64) (results []*Result, err error) {
	// run the document query
	docs, err := RunDocumentQuery(connection, q, from, size)
	if err != nil {
		return
	}

	// run the highlight queries
	highlights := make([]*HighlightResult, len(docs))
	for idx, doc := range docs {
		highlights[idx], err = doc.RunHighlightQuery(connection, q)
		if err != nil {
			return nil, err
		}
	}

	// and make proper results out of it
	results = make([]*Result, len(highlights))
	for idx, highlight := range highlights {
		results[idx], err = NewResult(highlight)
		if err != nil {
			return nil, err
		}
	}

	return
}

// NewResult generates a new result from a highlight result
func NewResult(highlight *HighlightResult) (result *Result, err error) {
	result = &Result{
		ElasticID: highlight.Document.ElasticID,
		Metadata:  highlight.Document.Element.Metadata,

		Score:    *highlight.Hit.Hit.Score,
		Snippets: highlight.Highlights,

		Math: highlight.Math,
	}
	return
}
