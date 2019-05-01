package query

// Query is a TemaSearch Query
type Query struct {
	Expressions []string // a list of MWS-expressions to search
	Text        string   // Text to search
}

// Kind represents the type of a query
type Kind int

const (
	// EmptyQueryKind is an empty query
	EmptyQueryKind Kind = iota
	// MWSQueryKind represents a MathWebSearch-only query
	MWSQueryKind
	// ElasticQueryKind represents an elastic-search only query
	ElasticQueryKind
	// TemaSearchQueryKind represents a dual (MWS, TemaSearch) query
	TemaSearchQueryKind
)

// Kind returns the type of this query
func (q *Query) Kind() Kind {
	hasMWS := q.NeedsMWS()
	hasElastic := q.NeedsElastic()

	if hasMWS && hasElastic {
		return TemaSearchQueryKind
	} else if hasMWS {
		return MWSQueryKind
	} else if hasElastic {
		return ElasticQueryKind
	}

	return EmptyQueryKind
}

// NeedsMWS checks if the query needs a MathWebSearch instance to resolve
func (q *Query) NeedsMWS() bool {
	return len(q.Expressions) != 0
}

// NeedsElastic checks if the query needs elasticsearch to resolve
func (q *Query) NeedsElastic() bool {
	return q.Text != ""
}

// MWSQuery turns this query into a MathWebSearch Query
func (q *Query) MWSQuery() *MWSQuery {
	return &MWSQuery{
		Expressions: q.Expressions,
		MwsIdsOnly:  q.Kind() == TemaSearchQueryKind,
	}
}

// ElasticQuery turns this query into an elastic query
func (q *Query) ElasticQuery(mwsIds []int64) *ElasticQuery {
	return &ElasticQuery{
		MathWebSearchIDs: mwsIds,
		Text:             q.Text,
	}
}
