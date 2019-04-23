package temasearch

import (
	"fmt"
	"sync/atomic"

	"github.com/MathWebSearch/mwsapi/mws"
	"github.com/MathWebSearch/mwsapi/tema/query"
	"github.com/MathWebSearch/mwsapi/utils"
)

// CountQuery counts the results of a query
func CountQuery(connection *Connection, q *Query) (count int64, err error) {
	// get the query type
	tp := q.Type()

	fmt.Printf("type = %#v query = %#v\n", tp, *q)

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

	// run at most size parallel operations
	size := connection.config.MWSPageSize
	err = utils.MaxParallel(int(connection.config.PoolSize), func(ch chan int64) {
		var i int64
		for i = 0; i <= outerTotal; i += size {
			ch <- i
		}
		close(ch)
	}, func(from int64) (e error) {
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

	return
}
