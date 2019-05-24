package elasticengine

import (
	"time"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/pkg/errors"

	"github.com/MathWebSearch/mwsapi/result"

	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
)

// runHighlightQuery runs a highlight query for a given element
func runHighlightQuery(conn *connection.ElasticConnection, q *query.ElasticQuery, hit *result.Hit) (took *time.Duration, err error) {
	// build the highlight query
	qq, h, err := q.RawHighlightQuery(hit)
	err = errors.Wrap(err, "q.RawHighlightQuery failed")
	if err != nil {
		return
	}

	// fetch the object and the highlights
	obj, took, err := elasticutils.FetchObject(conn.Client, conn.Config.HarvestIndex, conn.Config.HarvestType, qq, h)
	err = errors.Wrap(err, "elasticutils.FetchObject failed")
	if err == nil && obj == nil {
		err = errors.New("[runHighlightQuery] Can not find result")
	}

	if err != nil {
		return
	}

	// and build the highlight result
	err = hit.UnmarshalElasticHighlight(obj)
	err = errors.Wrap(err, "hit.UnmarshalElasticHighlight failed")
	return
}
