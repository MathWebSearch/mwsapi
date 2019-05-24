package temaengine

import (
	"fmt"
	"time"

	"github.com/MathWebSearch/mwsapi/engine/elasticengine"
	"github.com/MathWebSearch/mwsapi/engine/mwsengine"
	"github.com/pkg/errors"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/result"
)

// Run runs a query against elasticsearch + mws
func Run(conn *connection.TemaConnection, q *query.Query, from int64, size int64) (res *result.Result, err error) {
	res = &result.Result{From: from}

	// keep track of how long we took
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		res.Took = &duration
	}()

	// get the query type
	tp := q.Kind()

	fmt.Printf("type = %#v query = %#v\n", tp, *q) // TODO: Remove debug

	// empty query => return an empty result
	if tp == query.EmptyQueryKind {
		return
	}

	// mws query => run an mws query
	if tp == query.MWSQueryKind {
		err = runMWSQuery(conn, q, res, from, size)
		err = errors.Wrap(err, "runMWSQuery failed")
		return
	}

	// count elastic queries using the appropriate function
	if tp == query.ElasticQueryKind {
		err = runElasticQuery(conn, q, res, from, size)
		err = errors.Wrap(err, "runElasticQuery failed")
		return
	}

	// else run the multi-plexec query
	err = runTemaSearchQuery(conn, q, res, from, size)
	err = errors.Wrap(err, "runTemaSearchQuery failed")
	return
}

func runTemaSearchQuery(conn *connection.TemaConnection, q *query.Query, res *result.Result, from int64, size int64) (err error) {
	// build outer results
	qq := q.MWSQuery()
	outerTotal, err := mwsengine.Count(conn.MWS, qq)
	if err != nil || outerTotal == 0 {
		return
	}

	res.Kind = result.TemaSearchKind
	res.Hits = []*result.Hit{} // buffer for the inner hits

	// TODO: Allocate this into max size from the start, and then buffer
	outerfrom := int64(0) // current start index for outer queries
	innerfrom := from     // offset for the inner query

	nextinnerfrom := from    // the next inner 'from'
	hadinnerresults := false // did we take care of inner results

	res.Variables = []*result.QueryVariable{}
	setvariables := false

	maxHits := int(size)                                  // the maximum number of hits
	outerPageSize := conn.MWS.Config.MaxPageSize          // the page size for the outer pages
	innerPageSize := int(conn.Elastic.Config.MaxPageSize) // the page size for the inner pages

outer:
	for len(res.Hits) <= maxHits {
		// fetch the next outer page
		outer, err := mwsengine.Run(conn.MWS, qq, outerfrom, outerPageSize)
		err = errors.Wrap(err, "mwsengine.Run failed")
		outerfrom += outerPageSize
		if err != nil {
			return err
		}

		// set the variables
		if !setvariables {
			res.Variables = outer.Variables
			setvariables = true
		}

		// if we have no results, we have reached the end
		if len(outer.HitIDs) == 0 {
			break
		}

		// build the inner query
		qqq := q.ElasticQuery(outer.HitIDs)

		// prepare the next inner results
		hadinnerresults = false
		innerfrom = nextinnerfrom
	inner:
		for {
			// run the inner query
			// TODO: Re-do pagination
			inner, err := elasticengine.Run(conn.Elastic, qqq, innerfrom, conn.Elastic.Config.MaxPageSize)
			innerfrom += conn.Elastic.Config.MaxPageSize
			err = errors.Wrap(err, "elasticengine.Run failed")
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
			res.Hits = append(res.Hits, inner.Hits...)

			// if we have our number of results, break out of the outer loop
			if len(res.Hits) >= maxHits {
				break outer
			}

			// if we did not get a full page of results, or our result set is full, we break the cycle
			if len(inner.Hits) < innerPageSize {
				break inner
			}

		}

	}

	res.Hits = res.Hits[:maxHits]

	return
}

func runMWSQuery(conn *connection.TemaConnection, q *query.Query, res *result.Result, from int64, size int64) (err error) {
	// run the query
	result, err := mwsengine.Run(conn.MWS, q.MWSQuery(), from, size)
	err = errors.Wrap(err, "mwsengine.Run failed")
	if err != nil {
		return
	}

	// store the result
	*res = *result

	return
}

func runElasticQuery(conn *connection.TemaConnection, q *query.Query, res *result.Result, from int64, size int64) (err error) {
	// run the query
	result, err := elasticengine.Run(conn.Elastic, q.ElasticQuery(nil), from, size)
	err = errors.Wrap(err, "elasticengine.Run failed")
	if err != nil {
		return
	}

	// and copy over the elastic result
	*res = *result
	return
}
