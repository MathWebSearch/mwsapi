package elasticsync

import (
	"sync/atomic"

	"github.com/MathWebSearch/mwsapi/result"
	"github.com/MathWebSearch/mwsapi/utils/gogroup"
	"github.com/pkg/errors"

	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
	elastic "gopkg.in/olivere/elastic.v6"
)

// clearSegments clears untouched (old) segments from the index
func (proc *Process) clearSegments() (err error) {
	proc.Println(nil, "Removing old segments ...")

	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("touched", false))

	group := gogroup.NewWorkGroup(proc.conn.Config.PoolSize, true)

	old := elasticutils.FetchObjects(proc.conn.Client, proc.conn.Config.SegmentIndex, proc.conn.Config.SegmentType, q)
	for so := range old {
		func(so *elasticutils.Object) {
			job := gogroup.GroupJob(func(sync func(func())) error {
				atomic.AddInt64(&proc.stats.RemovedSegments, 1)
				return errors.Wrap(proc.clearSegment(sync, so), "proc.ClearSegment failed")
			})
			group.Add(&job)
		}(so)
	}

	err = group.Wait()
	if err == nil {
		proc.PrintlnOK(nil, "OK")
	} else {
		proc.PrintlnERROR(nil, "ERROR")
	}

	return
}

// clearSegment removes a single segment
func (proc *Process) clearSegment(sync func(func()), so *elasticutils.Object) (err error) {
	segment := &result.HarvestSegment{}
	err = so.Unpack(segment)
	err = errors.Wrap(err, "so.Unpack failed")
	if err != nil {
		return
	}
	proc.Printf(sync, "=> %s\n", segment.ID)

	// clear the harvests
	proc.Print(sync, "  Clearing harvests belonging to segment ... ")
	err = proc.clearSegmentHarvests(segment)
	err = errors.Wrap(err, "proc.clearSegmentHarvests failed")
	if err != nil {
		proc.PrintlnERROR(sync, "ERROR")
		return
	}
	proc.PrintlnOK(sync, "OK")

	// and remove segment itself
	proc.Print(sync, "  Removing segment ... ")
	err = so.Delete()
	err = errors.Wrap(err, "so.Delete failed")
	if err == nil {
		proc.PrintlnOK(sync, "OK")
	} else {
		proc.PrintlnERROR(sync, "ERROR")
	}

	return

}

// clearSegmentHarvests clears segments belonging to a harvest
func (proc *Process) clearSegmentHarvests(segment *result.HarvestSegment) (err error) {
	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("segment", segment.ID))

	err = elasticutils.DeleteBulk(proc.conn.Client, proc.conn.Config.HarvestIndex, proc.conn.Config.HarvestType, q)
	err = errors.Wrap(err, "elasticutils.DeleteBulk failed")
	return
}
