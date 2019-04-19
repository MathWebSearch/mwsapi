package query

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MathWebSearch/mwsapi/elasticutils"
	"github.com/MathWebSearch/mwsapi/tema"

	"gopkg.in/olivere/elastic.v6"
)

// DocumentResult represents a result of the Document phase
type DocumentResult struct {
	Total int64 `json:"total"` // the total number of results
	From  int64 `json:"from"`  // result number this page starts at
	Size  int64 `json:"size"`  // (maximum) number of results in this page

	Took *time.Duration `json:"took"` // the amount of time the query took to execute, including network latency to elasticsearch

	Hits []*DocumentHit `json:"hits"` // the current page of results
}

// DocumentHit represents the result of a documentquery
type DocumentHit struct {
	ID      string               `json:"id"`     // id of the element being returned
	Element *tema.HarvestElement `json:"source"` // source of the found element

	Math []*FormulaeInfo `json:"math"` // the list of math elements within this hit
}

// FormulaeInfo represents a single math excert within an element
type FormulaeInfo struct {
	MathID string `json:"id"`    // id of this element
	XPath  string `json:"xpath"` // path of this element
}

// RealMathID return the real math id
func (info *FormulaeInfo) RealMathID() string {
	if _, err := strconv.Atoi(info.MathID); err == nil {
		return "math" + info.MathID
	}
	return info.MathID
}

// RunDocumentQuery runs the document query phase of a query
func RunDocumentQuery(connection *tema.Connection, query *Query, from int64, size int64) (result *DocumentResult, err error) {

	// measure time for this query
	start := time.Now()
	defer func() {
		took := time.Since(start)
		result.Took = &took
	}()

	result = &DocumentResult{
		From: from,
		Size: size,
	}

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

	// prepare result objejct
	result.Total = page.Total
	result.Hits = make([]*DocumentHit, len(page.Hits))

	// and make the new hits
	for i, hit := range page.Hits {
		result.Hits[i], err = NewDocumentHit(hit)
		if err != nil {
			return nil, err
		}
	}

	return
}

// CountDocumentQuery counts all obects subject to a document query
func CountDocumentQuery(connection *tema.Connection, query *Query) (count int64, err error) {
	// make the document query
	q, err := query.asDocumentQuery()
	if err != nil {
		return
	}

	// and run the count utility
	return elasticutils.Count(connection.Client, connection.Config.HarvestIndex, connection.Config.HarvestType, q)
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

// NewDocumentHit generates a hit from an elasticsearch object
func NewDocumentHit(obj *elasticutils.Object) (result *DocumentHit, err error) {
	result = &DocumentHit{
		ID: obj.GetID(),
	}

	// store the source element
	var raw tema.HarvestElement
	err = obj.Unpack(&raw)
	if err != nil {
		return
	}
	result.Element = &raw

	// create the math elements
	// we do not a-priori know the size of it
	result.Math = []*FormulaeInfo{}

	for _, mwsid := range result.Element.MWSNumbers {
		// load the data
		data, ok := result.Element.MWSPaths[mwsid]
		if !ok {
			return nil, fmt.Errorf("Result %q missing path info for %d", result.ID, mwsid)
		}

		// and iterate over it
		for key, value := range data {
			result.Math = append(result.Math, &FormulaeInfo{
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
