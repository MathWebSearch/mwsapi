package temasearch

import (
	"fmt"
	"time"

	"github.com/MathWebSearch/mwsapi/mws"
	"github.com/MathWebSearch/mwsapi/tema/query"
)

// RunQuery runs a query
func RunQuery(connection *Connection, q *Query, from int64, size int64) (res *Result, err error) {
	res = &Result{From: from}

	// keep track of how long we took
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		res.Took = &duration
	}()

	// get the query type
	tp := q.Type()

	fmt.Printf("type = %#v query = %#v\n", tp, *q)

	// empty query => return an empty result
	if tp == EmptyQuery {
		return
	}

	// mws query => run an mws query
	if tp == MWSQuery {
		err = runMWSQuery(connection, q, res, from, size)
		return
	}

	// count elastic queries using the appropriate function
	if tp == ElasticQuery {
		err = runElasticQuery(connection, q, res, from, size)
		return
	}

	// else run the multi-plexec query
	err = runTemaSearchQuery(connection, q, res, from, size)
	return
}

func runTemaSearchQuery(connection *Connection, q *Query, res *Result, from int64, size int64) (err error) {
	return
}

func runMWSQuery(connection *Connection, q *Query, res *Result, from int64, size int64) (err error) {
	// run the query
	result, err := mws.RunQuery(connection.mws, q.asMWSQuery(), from, size)
	if err != nil {
		return
	}

	// and copy over the mathwebsearch result
	res.fromMathWebSearch(result)
	return
}

func runElasticQuery(connection *Connection, q *Query, res *Result, from int64, size int64) (err error) {
	// run the query
	result, err := query.RunQuery(connection.tema, q.asElasticQuery(nil), from, size)
	if err != nil {
		return
	}

	// and copy over the elastic result
	res.fromElastic(result)
	return
}
