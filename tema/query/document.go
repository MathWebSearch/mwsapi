package query

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MathWebSearch/mwsapi/elasticutils"
	"github.com/MathWebSearch/mwsapi/tema"

	"gopkg.in/olivere/elastic.v6"
)

// DocumentResult represents the result of a documentquery
type DocumentResult struct {
	ElasticID string
	Element   *tema.HarvestElement

	Math []*MathDocumentInfo
}

// MathDocumentInfo represents a single math excert within an element
type MathDocumentInfo struct {
	MathID string
	XPath  string
}

// RealMathID return the real math id
func (info *MathDocumentInfo) RealMathID() string {
	if _, err := strconv.Atoi(info.MathID); err == nil {
		return "math" + info.MathID
	}
	return info.MathID
}

// RunDocumentQuery runs the document query phase of a query
func RunDocumentQuery(connection *tema.Connection, query *Query, from int64, size int64) (results []*DocumentResult, err error) {
	// make the document query
	q, err := query.asDocumentQuery()
	if err != nil {
		return
	}

	// grab the results
	page, err := elasticutils.FetchObjectsPage(connection.Client, connection.Config.HarvestIndex, connection.Config.HarvestType, q, nil, from, size)
	if err != nil {
		return
	}

	// make a document result slice
	results = make([]*DocumentResult, len(page.Hits))

	for i, hit := range page.Hits {
		results[i], err = NewHit(hit)
		if err != nil {
			return nil, err
		}
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
	if len(query.MathWebSearchIDs) > 0 {
		// need to convert []int64 to []interface{}
		ids := make([]interface{}, len(query.MathWebSearchIDs))
		for i, v := range query.MathWebSearchIDs {
			ids[i] = v
		}

		formulae := elastic.NewTermsQuery("mws_ids", ids...)
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
	result.Element = &raw

	for _, mwsid := range result.Element.MWSNumbers {
		// load the data
		data, ok := result.Element.MWSPaths[mwsid]
		if !ok {
			return nil, fmt.Errorf("Result %q missing path info for %d", result.ElasticID, mwsid)
		}

		// and iterate over it
		for key, value := range data {
			result.Math = append(result.Math, &MathDocumentInfo{
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
