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
	// build outer results
	qq := q.asMWSQuery()
	outerTotal, err := mws.CountQuery(connection.mws, qq)
	if err != nil || outerTotal == 0 {
		return
	}

	hits := []*query.Hit{} // buffer for the inner hits

	outerfrom := int64(0) // current start index for outer queries
	innerfrom := from     // offset for the inner query

	nextinnerfrom := from    // the next inner 'from'
	hadinnerresults := false // did we take care of inner results

	variables := []*mws.QueryVariable{}
	setvariables := false

	maxHits := int(size)                                 // the maximum number of hits
	outerPageSize := connection.config.MWSPageSize       // the page size for the outer pages
	innerPageSize := int(connection.config.TemaPageSize) // the page size for the inner pages

outer:
	for len(hits) <= maxHits {
		// fetch the next outer page
		outer, err := mws.RunQuery(connection.mws, qq, outerfrom, outerPageSize)
		outerfrom += outerPageSize
		if err != nil {
			return err
		}

		// set the variables
		if !setvariables {
			variables = outer.Variables
			setvariables = true
		}

		// if we have no results, we have reached the end
		if len(outer.MathWebSearchIDs) == 0 {
			break
		}

		// build the inner query
		qqq := q.asElasticQuery(outer.MathWebSearchIDs)

		// prepare the next inner results
		hadinnerresults = false
		innerfrom = nextinnerfrom
	inner:
		for {
			// run the inner query
			inner, err := query.RunQuery(connection.tema, qqq, innerfrom, connection.config.TemaPageSize)
			innerfrom += connection.config.TemaPageSize
			if err != nil {
				return err
			}

			// prepare the next inner from clause
			// which can skip these results
			if !hadinnerresults {
				nextinnerfrom = nextinnerfrom - inner.Total
				hadinnerresults = true
			}

			// append all the hits and prepare for the next page
			hits = append(hits, inner.Hits...)

			// if we have our number of results, break out of the outer loop
			if len(hits) >= maxHits {
				break outer
			}

			// if we did not get a full page of results, or our result set is full, we break the cycle
			if len(inner.Hits) < innerPageSize {
				break inner
			}

		}

	}

	// slice the result set to never return more than what was requested
	res.fromTema(variables, hits[:maxHits])

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
