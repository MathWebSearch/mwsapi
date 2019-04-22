package mws

// Query Represents a MathWebSearch Query
type Query struct {
	Expressions []string // MathWebSearch Expressions to query for
	MwsIdsOnly  bool     // if set to true, use method "mws_ids", else "json"
}

// RunQuery runs a query
func RunQuery(connection *Connection, q *Query, from int64, size int64) (res *Result, err error) {
	return connection.runRawQuery(q.asNewRawQuery(from, size))
}

// CountQuery counts the results in a single query
func CountQuery(connection *Connection, q *Query) (size int64, err error) {

	// run a query returning 0 results, but still counting
	res, err := RunQuery(connection, q, 0, 0)
	if err != nil {
		return
	}

	// and count
	size = res.Total
	return
}

func (q *Query) asNewRawQuery(from int64, size int64) *RawQuery {
	// make the expressions
	exprs := make([]*Expression, len(q.Expressions))
	for i, expr := range q.Expressions {
		exprs[i] = &Expression{
			Term: expr,
		}
	}

	var format string
	if q.MwsIdsOnly {
		format = "mws-ids"
	} else {
		format = "json"
	}

	// and make the new raw query
	return &RawQuery{
		From: from,
		Size: size,

		ReturnTotal:  true,
		OutputFormat: format,

		Expressions: exprs,
	}
}
