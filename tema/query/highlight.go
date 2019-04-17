package query

import (
	"errors"
	"fmt"

	"github.com/MathWebSearch/mwsapi/elasticutils"
	"github.com/MathWebSearch/mwsapi/tema"
	"gopkg.in/olivere/elastic.v6"
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

// RunHighlightQuery runs a highlight query for a given result
func (res *DocumentResult) RunHighlightQuery(connection *tema.Connection, query *Query) (result *HighlightResult, err error) {
	// build the highlight query
	q, h, err := query.asHighlightQuery(res)
	if err != nil {
		return
	}

	// fetch the object and the highlights
	obj, err := elasticutils.FetchObject(connection.Client, connection.Config.HarvestIndex, connection.Config.HarvestType, q, h)
	if err == nil && obj == nil {
		err = errors.New("Can not find result")
	}

	if err != nil {
		return nil, err
	}

	result, err = NewHighlightResult(res, obj)
	return
}

func (query *Query) asHighlightQuery(res *DocumentResult) (elastic.Query, *elastic.Highlight, error) {
	q := elastic.NewBoolQuery().Must(elastic.NewIdsQuery().Ids(res.ElasticID))
	nonEmptyQuery := false

	// text highlights first
	if len(query.Text) > 0 {
		text := elastic.NewMatchQuery("text", query.Text)
		q = q.Must(text)
		nonEmptyQuery = true
	}

	// formulae highlights next
	for _, math := range res.Math {
		matcher := elastic.NewMatchQuery("text", math.RealMathID()).MinimumShouldMatch("2").Operator("or")
		q = q.Must(matcher)
	}

	if !nonEmptyQuery {
		return nil, nil, errors.New("Query had neither text nor mws_ids")
	}

	// build the highlight itself
	h := elastic.NewHighlight().Fields(elastic.NewHighlighterField("text")).PreTags("<span class=\"tema-highlight\">").PostTags("</span>")

	// and return the query itself
	return q, h, nil
}

// NewHighlightResult makes a new highlight result
func NewHighlightResult(doc *DocumentResult, obj *elasticutils.Object) (res *HighlightResult, err error) {
	if obj.Hit == nil || obj.Hit.Highlight == nil {
		return nil, errors.New("No highlights returned")
	}

	res = &HighlightResult{
		Document: doc,
		Hit:      obj,
	}

	// load the highlights
	var ok bool
	res.Highlights, ok = (*obj.Hit.Highlight)["text"]
	if !ok {
		return nil, errors.New("No highlights returned")
	}

	// map() over doc.Math
	res.Math = make([]*ReplacedMath, len(doc.Math))
	for i, math := range doc.Math {
		res.Math[i] = &ReplacedMath{
			ID:    math.RealMathID(),
			XPath: math.XPath,
		}

		var ok bool
		res.Math[i].Source, ok = doc.Element.MathSource[math.MathID]
		if !ok {
			return nil, fmt.Errorf("Result %s has no source info for %s", doc.ElasticID, math.MathID)
		}
	}

	return
}
