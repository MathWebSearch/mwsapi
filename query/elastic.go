package query

import (
	"errors"

	"github.com/MathWebSearch/mwsapi/result"
	elastic "gopkg.in/olivere/elastic.v6"
)

// ElasticQuery represents an elasticsearch query
type ElasticQuery struct {
	MathWebSearchIDs []int64 // list of possible formula ids
	Text             string  // Textual content
}

// RawElasticQuery represents a raw query sent to elastic query
type RawElasticQuery elastic.Query

// RawDocumentQuery turns this query into a RawDocumentQuery
func (q *ElasticQuery) RawDocumentQuery() (RawElasticQuery, error) {
	query := elastic.NewBoolQuery()

	nonEmptyQuery := false

	// if we have some formulae
	if len(q.Text) > 0 {
		text := elastic.NewMatchQuery("text", q.Text).MinimumShouldMatch("2").Operator("or")
		query = query.Must(text)
		nonEmptyQuery = true
	}

	// and return the formula id
	if len(q.MathWebSearchIDs) > 0 {
		// need to convert []int64 to []interface{}
		ids := make([]interface{}, len(q.MathWebSearchIDs))
		for i, v := range q.MathWebSearchIDs {
			ids[i] = v
		}

		formulae := elastic.NewTermsQuery("mws_ids", ids...)
		query = query.Must(formulae)
		nonEmptyQuery = true
	}

	if !nonEmptyQuery {
		return nil, errors.New("[ElasticQuery.RawDocumentQuery] Query had neither text nor mws_ids")
	}

	// and return the query itself
	return query, nil
}

// RawHighlightQuery turns this query into a HighlightQuery
func (q *ElasticQuery) RawHighlightQuery(res *result.Hit) (RawElasticQuery, *elastic.Highlight, error) {
	q2 := elastic.NewBoolQuery().Must(elastic.NewIdsQuery().Ids(res.ID))
	nonEmptyQuery := false

	// text highlights first
	if len(q.Text) > 0 {
		text := elastic.NewMatchQuery("text", q.Text)
		q2 = q2.Must(text)
		nonEmptyQuery = true
	}

	// formulae highlights next
	for _, math := range res.Math {
		matcher := elastic.NewMatchQuery("text", math.RealMathID()).MinimumShouldMatch("2").Operator("or")
		q2 = q2.Must(matcher)
		nonEmptyQuery = true
	}

	if !nonEmptyQuery {
		return nil, nil, errors.New("[ElasticQuery.RawHighlightQuery] Neither text nor ids provided")
	}

	// build the highlight itself
	h := elastic.NewHighlight().Fields(elastic.NewHighlighterField("text")).PreTags("<span class=\"tema-highlight\">").PostTags("</span>")

	// and return the query itself
	return q2, h, nil
}
