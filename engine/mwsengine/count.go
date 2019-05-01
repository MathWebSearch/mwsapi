package mwsengine

import (
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/query"
)

// Count counts MathWebSearch Query results
func Count(conn *connection.MWSConnection, query *query.MWSQuery) (count int64, err error) {
	res, err := Run(conn, query, 0, 0)
	if err != nil {
		return
	}

	count = res.Total
	return
}
