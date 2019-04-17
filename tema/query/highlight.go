package query

import (
	"errors"

	"github.com/MathWebSearch/mwsapi/elasticutils"
	"github.com/MathWebSearch/mwsapi/tema"
	"github.com/olivere/elastic"
)

// HighlightResult represents the
type HighlightResult struct {
	// the document and the math this result came from
	Document *DocumentResult

	// the replaced math
	Math []*ReplacedMath

	// The element that was being 'hit'
	Hit *elasticutils.Object

	// the highlights themselves
	Highlights []string
}

// ReplacedMath represents a single replaced MathExcept
type ReplacedMath struct {
	Source string
	ID     string
	XPath  string
}

// RunHighlightQuery runs a highlight query for a given result
func (res *DocumentResult) RunHighlightQuery(connection *tema.Connection, query *Query) (result *HighlightResult, err error) {
	q, err := query.asHighlightQuery(res)

	// fetch the object and the highlights
	obj, err := elasticutils.FetchObject(connection.Client, connection.Config.HarvestIndex, connection.Config.HarvestType, q)
	if err == nil && result == nil {
		err = errors.New("Unable to highlight result")
	}

	if err != nil {
		return nil, err
	}

	result, err = NewHighlightResult(res, obj)
	return
}

func (query *Query) asHighlightQuery(res *DocumentResult) (elastic.Query, error) {
	q := elastic.NewBoolQuery()
	nonEmptyQuery := false

	// if we have some formulae
	if len(query.Text) > 0 {
		text := elastic.NewMatchQuery("text", query.Text).MinimumShouldMatch("2").Operator("or")
		q = q.Must(text)
		nonEmptyQuery = true
	}

	// and return the formula id
	if len(query.FormulaID) > 0 {
		formulae := elastic.NewTermQuery("mws_ids", query.FormulaID)
		q = q.Must(formulae)
		nonEmptyQuery = true
	}

	if !nonEmptyQuery {
		return nil, errors.New("Query had neither text nor mws_ids")
	}

	// and return the query itself
	return q, nil
}

// NewHighlightResult makes a new highlight result
func NewHighlightResult(doc *DocumentResult, obj *elasticutils.Object) (res *HighlightResult, err error) {
	err = errors.New("Not implemented")
	return
}
