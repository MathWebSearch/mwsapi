package temasearch

import (
	"sync/atomic"

	"github.com/MathWebSearch/mwsapi/utils/gogroup"

	"github.com/MathWebSearch/mwsapi/mws"
	"github.com/MathWebSearch/mwsapi/tema/query"
)

// CountQuery counts the results of a query
func CountQuery(connection *Connection, q *Query) (count int64, err error) {
	// get the query type
	tp := q.Type()

	// empty queries have 0 results
	if tp == EmptyQuery {
		return 0, nil
	}

	// count mws queries using the appropriate function
	if tp == MWSQuery {
		return countMWSQuery(connection, q)
	}

	// count elastic queries using the appropriate function
	if tp == ElasticQuery {
		return countElasticQuery(connection, q, nil)
	}

	// returns
	return countTemaSearchQuery(connection, q)
}

func countMWSQuery(connection *Connection, q *Query) (int64, error) {
	return mws.CountQuery(connection.mws, q.asMWSQuery())
}

func countElasticQuery(connection *Connection, q *Query, mwsIds []int64) (int64, error) {
	return query.CountQuery(connection.tema, q.asElasticQuery(mwsIds))
}

func countTemaSearchQuery(connection *Connection, q *Query) (count int64, err error) {
	// build the outer query
	qq := q.asMWSQuery()

	// query the total number of outer results
	outerTotal, err := mws.CountQuery(connection.mws, qq)
	if err != nil || outerTotal == 0 {
		return
	}

	// create a group for at most (PoolSize) parallel operations
	group := gogroup.NewWorkGroup(int(connection.config.PoolSize), false)

	// and add the jobs
	size := connection.config.MWSPageSize
	for i := int64(0); i <= outerTotal; i += size {
		(func(from int64) {
			job := gogroup.GroupJob(func(sync func(func())) (e error) {
				// run the outer query and exit if it has an empty result
				outer, e := mws.RunQuery(connection.mws, qq, from, connection.config.MWSPageSize)
				if e != nil || len(outer.MathWebSearchIDs) == 0 {
					return
				}

				// get the total number of inner results
				innertotal, e := query.CountQuery(connection.tema, q.asElasticQuery(outer.MathWebSearchIDs))
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
