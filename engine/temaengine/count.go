package temaengine

import (
	"sync/atomic"

	"github.com/MathWebSearch/mwsapi/engine/elasticengine"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/engine/mwsengine"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/utils/gogroup"
)

// Count counts the results of a tema query
func Count(conn *connection.TemaConnection, q *query.Query) (count int64, err error) {
	// get the query type
	tp := q.Kind()

	// empty queries have 0 results
	if tp == query.EmptyQueryKind {
		return 0, nil
	}

	// count mws queries using the appropriate function
	if tp == query.MWSQueryKind {
		return countMWSQuery(conn, q)
	}

	// count elastic queries using the appropriate function
	if tp == query.ElasticQueryKind {
		return countElasticQuery(conn, q, nil)
	}

	// returns
	return countTemaSearchQuery(conn, q)
}

func countMWSQuery(conn *connection.TemaConnection, q *query.Query) (int64, error) {
	return mwsengine.Count(conn.MWS, q.MWSQuery())
}

func countElasticQuery(conn *connection.TemaConnection, q *query.Query, mwsIds []int64) (int64, error) {
	return elasticengine.Count(conn.Elastic, q.ElasticQuery(mwsIds))
}

func countTemaSearchQuery(conn *connection.TemaConnection, q *query.Query) (count int64, err error) {
	// build the outer query
	qq := q.MWSQuery()

	// query the total number of outer results
	outerTotal, err := mwsengine.Count(conn.MWS, qq)
	if err != nil || outerTotal == 0 {
		return
	}

	// create a group for at most (PoolSize) parallel operations
	group := gogroup.NewWorkGroup(conn.MWS.Config.PoolSize, false)

	// and add the jobs
	// TODO: Generalize this pagination
	size := int64(conn.MWS.Config.PoolSize)
	for i := int64(0); i <= outerTotal; i += size {
		(func(from int64) {
			job := gogroup.GroupJob(func(sync func(func())) (e error) {
				// run the outer query and exit if it has an empty result
				outer, e := mwsengine.Run(conn.MWS, qq, from, conn.MWS.Config.MaxPageSize)
				if e != nil || len(outer.HitIDs) == 0 {
					return
				}

				// get the total number of inner results
				innertotal, e := elasticengine.Count(conn.Elastic, q.ElasticQuery(outer.HitIDs))
				if e != nil {
					return
				}

				// add them and return
				atomic.AddInt64(&count, innertotal)
				return
			})
			group.Add(&job)
		})(i)
	}

	// then wait
	err = group.Wait()
	return
}
