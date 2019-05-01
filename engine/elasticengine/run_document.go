package elasticengine

import (
	"time"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/result"
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
)

// RunDocument runs the document query
func RunDocument(conn *connection.ElasticConnection, q *query.ElasticQuery, from int64, size int64) (res *result.Result, err error) {
	// measure time for this query
	start := time.Now()
	defer func() {
		if res != nil {
			took := time.Since(start)
			res.Took = &took
		}
	}()

	res = &result.Result{
		From: from,
	}

	// make the document query
	qq, err := q.RawDocumentQuery()
	if err != nil {
		return
	}

	// grab the results
	// TODO: Paralellize this with appropriate page size
	page, err := elasticutils.FetchObjectsPage(conn.Client, conn.Config.HarvestIndex, conn.Config.HarvestType, qq, nil, from, size)
	if err != nil {
		return
	}

	err = newDocumentResult(res, page)

	return
}
