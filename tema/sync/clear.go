package sync

import (
	"sync/atomic"

	"github.com/MathWebSearch/mwsapi/elasticutils"
	"github.com/MathWebSearch/mwsapi/tema"
	"github.com/MathWebSearch/mwsapi/utils"
	"gopkg.in/olivere/elastic.v6"
)

// clearSegments clears untouched (old) segments from the index
func (proc *Process) clearSegments() (err error) {
	proc.println(nil, "Removing old segments ...")

	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("touched", false))

	group := utils.NewAsyncGroup()

	old := elasticutils.FetchObjects(proc.connection.Client, proc.connection.Config.SegmentIndex, proc.connection.Config.SegmentType, q)
	for so := range old {
		func(so *elasticutils.Object) {
			group.Add(func(sync func(func())) error {
				atomic.AddInt64(&proc.stats.RemovedSegments, 1)
				return proc.clearSegment(sync, so)
			})
		}(so)
	}

	err = group.Wait()
	if err == nil {
		proc.printlnOK(nil, "OK")
	} else {
		proc.printlnERROR(nil, "ERROR")
	}

	return
}

// clearSegment removes a single segment
func (proc *Process) clearSegment(sync func(func()), so *elasticutils.Object) (err error) {
	segment := &tema.Segment{}
	err = so.Unpack(segment)
	if err != nil {
		return
	}
	proc.printf(sync, "=> %s\n", segment.ID)

	// clear the harvests
	proc.print(sync, "  Clearing harvests belonging to segment ... ")
	err = proc.clearSegmentHarvests(segment)
	if err != nil {
		proc.printlnERROR(sync, "ERROR")
		return
	}
	proc.printlnOK(sync, "OK")

	// and remove segment itself
	proc.print(sync, "  Removing segment ... ")
	err = so.Delete()
	if err == nil {
		proc.printlnOK(sync, "OK")
	} else {
		proc.printlnERROR(sync, "ERROR")
	}

	return

}

// clearSegmentHarvests clears segments belonging to a harvest
func (proc *Process) clearSegmentHarvests(segment *tema.Segment) error {
	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("segment", segment.ID))

	return elasticutils.DeleteBulk(proc.connection.Client, proc.connection.Config.HarvestIndex, proc.connection.Config.HarvestType, q)
}
