package elasticengine

import (
	"errors"
	"time"

	"github.com/MathWebSearch/mwsapi/connection"

	"github.com/MathWebSearch/mwsapi/result"

	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
)

// runHighlightQuery runs a highlight query for a given element
func runHighlightQuery(conn *connection.TemaConnection, q *query.ElasticQuery, hit *result.Hit) (res *result.Hit, took *time.Duration, err error) {
	// build the highlight query
	qq, h, err := q.RawHighlightQuery(hit)
	if err != nil {
		return
	}

	// fetch the object and the highlights
	obj, took, err := elasticutils.FetchObject(conn.Client, conn.Config.HarvestIndex, conn.Config.HarvestType, qq, h)
	if err == nil && obj == nil {
		err = errors.New("Can not find result")
	}

	if err != nil {
		return
	}

	// and build the highlight result
	err = newHighlightHit(hit, obj)
	if err != nil {
		return
	}

	// and take over the result
	res = hit
	return
}
