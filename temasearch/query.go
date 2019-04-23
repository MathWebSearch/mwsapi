package temasearch

import (
	"github.com/MathWebSearch/mwsapi/mws"
	"github.com/MathWebSearch/mwsapi/tema/query"
)

// Query is a TemaSearch Query
type Query struct {
	Expressions []string // a list of expressions to search
	Text        string   // some MathWebSearch Text to search for
}

// QueryType represents the type of a query
type QueryType int

const (
	// EmptyQuery is an empty query
	EmptyQuery QueryType = iota
	// MWSQuery represents a MathWebSearch-only query
	MWSQuery
	// ElasticQuery represents an elastic-search only query
	ElasticQuery
	// TemaSearchQuery represents a dual (MWS, TemaSearch) query
	TemaSearchQuery
)

// Type returns the type of this query
func (q *Query) Type() QueryType {
	hasMWS := q.NeedsMWS()
	hasElastic := q.NeedsElastic()

	if hasMWS && hasElastic {
		return TemaSearchQuery
	} else if hasMWS {
		return MWSQuery
	} else if hasElastic {
		return ElasticQuery
	}

	return EmptyQuery
}

// NeedsMWS checks if the query needs a MathWebSearch instance to resolve
func (q *Query) NeedsMWS() bool {
	return len(q.Expressions) != 0
}

// NeedsElastic checks if the query needs elasticsearch to resolve
func (q *Query) NeedsElastic() bool {
	return q.Text != ""
}

// asMWSQuery turns this query into a MathWebSearch Query
func (q *Query) asMWSQuery() *mws.Query {
	return &mws.Query{
		Expressions: q.Expressions,
		MwsIdsOnly:  q.Type() == TemaSearchQuery,
	}
}

// asElasticQuery turns this query into an elastic query
func (q *Query) asElasticQuery(mwsIds []int64) *query.Query {
	return &query.Query{
		MathWebSearchIDs: mwsIds,
		Text:             q.Text,
	}
}
