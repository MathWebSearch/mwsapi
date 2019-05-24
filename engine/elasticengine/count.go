package elasticengine

import (
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
	"github.com/pkg/errors"
)

// Count counts the number of results for an elastic query
func Count(conn *connection.ElasticConnection, q *query.ElasticQuery) (count int64, err error) {
	qq, err := q.RawDocumentQuery()
	err = errors.Wrap(err, "q.RawDocumentQuery failed")
	if err != nil {
		return
	}

	count, err = elasticutils.Count(conn.Client, conn.Config.HarvestIndex, conn.Config.HarvestType, qq)
	err = errors.Wrap(err, "elasticutils.Count failed")
	return
}
