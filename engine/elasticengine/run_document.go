package elasticengine

import (
	"time"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/result"
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
	"github.com/pkg/errors"
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
	err = errors.Wrap(err, "a.RawDocumentQuery failed")
	if err != nil {
		return
	}

	// grab the results
	// TODO: Paralellize this with appropriate page size
	page, err := elasticutils.FetchObjectsPage(conn.Client, conn.Config.HarvestIndex, conn.Config.HarvestType, qq, nil, from, size)
	err = errors.Wrap(err, "elasticutils.FetchObjectsPage failed")
	if err != nil {
		return
	}

	// and un-marshal the results
	err = res.UnmarshalElastic(page)
	err = errors.Wrap(err, "res.UnmarshalElastic failed")
	return
}
