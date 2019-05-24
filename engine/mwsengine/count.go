package mwsengine

import (
	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/pkg/errors"
)

// Count counts MathWebSearch Query results
func Count(conn *connection.MWSConnection, query *query.MWSQuery) (count int64, err error) {
	res, err := Run(conn, query, 0, 0)
	err = errors.Wrap(err, "Run failed")
	if err != nil {
		return
	}

	count = res.Total
	return
}
