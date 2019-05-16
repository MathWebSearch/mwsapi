package elasticengine

import (
	"sync"
	"time"

	"github.com/MathWebSearch/mwsapi/connection"
	"github.com/MathWebSearch/mwsapi/query"
	"github.com/MathWebSearch/mwsapi/result"
	"github.com/MathWebSearch/mwsapi/utils/gogroup"
)

// Run runs an elasticsearch query
func Run(conn *connection.ElasticConnection, q *query.ElasticQuery, from int64, size int64) (res *result.Result, err error) {
	// measure time for this query
	start := time.Now()
	defer func() {
		if res != nil {
			took := time.Since(start)
			res.Took = &took
		}
	}()

	res = &result.Result{
		From: from,
		Size: size,
	}

	// run the document query
	res, err = RunDocument(conn, q, from, size)
	if err != nil {
		return
	}

	// prepare running the highlight query in parallel
	group := gogroup.NewWorkGroup(conn.Config.PoolSize, false)

	tookMutex := &sync.Mutex{}
	tookTotal := time.Duration(0)

	for idx, doc := range res.Hits {
		func(idx int, doc *result.Hit) {
			job := gogroup.GroupJob(func(_ func(func())) (err error) {
				var took *time.Duration

				took, err = runHighlightQuery(conn, q, doc)

				// add the total time taken
				if err == nil {
					tookMutex.Lock()
					tookTotal += *took
					tookMutex.Unlock()
				}

				return
			})
			group.Add(&job)
		}(idx, doc)
	}

	// wait for them all to come back
	err = group.Wait()
	if err != nil {
		return
	}

	// add the time taken
	res.Kind = result.ElasticKind
	res.TookComponents["highlight"] = &tookTotal

	return
}
