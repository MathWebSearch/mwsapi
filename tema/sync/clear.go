package sync

import (
	"github.com/MathWebSearch/mwsapi/elasticutils"
	"github.com/MathWebSearch/mwsapi/tema"
	"gopkg.in/olivere/elastic.v6"
)

// clearSegments clears untouched (old) segments from the index
func (proc *Process) clearSegments() (err error) {
	proc.println("Removing old segments ...")

	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("touched", false))

	old := elasticutils.FetchObjects(proc.connection.Client, proc.connection.Config.SegmentIndex, proc.connection.Config.SegmentType, q)
	for so := range old {
		proc.stats.RemovedSegments++
		e := proc.clearSegment(so)
		if e != nil {
			err = e
		}
	}

	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return
}

// clearSegment removes a single segment
func (proc *Process) clearSegment(so *elasticutils.Object) (err error) {
	segment := &tema.Segment{}
	err = so.Unpack(segment)
	if err != nil {
		return
	}
	proc.printf("=> %s\n", segment.ID)

	// clear the harvests
	proc.print("  Clearing harvests belonging to segment ... ")
	err = proc.clearSegmentHarvests(segment)
	if err != nil {
		proc.printlnERROR("ERROR")
		return
	}
	proc.printlnOK("OK")

	// and remove segment itself
	proc.print("  Removing segment ... ")
	err = so.Delete()
	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return

}

// clearSegmentHarvests clears segments belonging to a harvest
func (proc *Process) clearSegmentHarvests(segment *tema.Segment) error {
	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("segment", segment.ID))

	return elasticutils.DeleteBulk(proc.connection.Client, proc.connection.Config.HarvestIndex, proc.connection.Config.HarvestType, q)
}
