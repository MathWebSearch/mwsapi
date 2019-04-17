package query

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MathWebSearch/mwsapi/elasticutils"
	"github.com/MathWebSearch/mwsapi/tema"

	"github.com/olivere/elastic"
)

// DocumentResult represents the result of a documentquery
type DocumentResult struct {
	ElasticID string

	Math []*MathExcept
}

// MathExcept represents a single math excert within an element
type MathExcept struct {
	MathID string
	XPath  string
}

// RunDocumentQuery runs the document query phase of a query
func RunDocumentQuery(connection *tema.Connection, query *Query, from int, size int) (results []*DocumentResult, err error) {
	// make the document query
	q, err := query.asDocumentQuery()
	if err != nil {
		return
	}

	// grab the results
	page, err := elasticutils.FetchObjectsPage(connection.Client, connection.Config.HarvestIndex, connection.Config.HarvestType, q, from, size)
	if err != nil {
		return
	}

	for _, hit := range page.Hits {
		doc, err := NewHit(hit)
		if err != nil {
			return nil, err
		}

		// and append the result
		results = append(results, doc)
	}

	return
}

func (query *Query) asDocumentQuery() (elastic.Query, error) {
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

// NewHit generates a hit from an elasticsearch object
func NewHit(obj *elasticutils.Object) (result *DocumentResult, err error) {
	result = &DocumentResult{
		ElasticID: obj.GetID(),
	}

	var raw tema.HarvestElement
	err = obj.Unpack(&raw)
	if err != nil {
		return
	}

	for _, mwsid := range raw.MWSNumbers {
		// load the data
		data, ok := raw.MWSData[mwsid]
		if !ok {
			return nil, fmt.Errorf("Result %q missing data for %d", result.ElasticID, mwsid)
		}

		// and iterate over it
		for key, value := range data {
			result.Math = append(result.Math, &MathExcept{
				MathID: simplifyMathID(key),
				XPath:  value.XPath,
			})
		}
	}

	return
}

func simplifyMathID(id string) string {
	parts := strings.Split(id, "#")
	return parts[len(parts)-1]
}
