package elasticsync

import (
	"github.com/MathWebSearch/mwsapi/result"
	elastic "gopkg.in/olivere/elastic.v6"

	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
)

// checkSegmentIndex checks the segment index for a given segment
func (proc *Process) checkSegmentIndex(segment string) (segobj *result.HarvestSegment, obj *elasticutils.Object, created bool, err error) {
	// the query
	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("segment", segment))

	// the fields
	segFields := result.HarvestSegment{
		ID:      segment,
		Hash:    "",
		Touched: true,
	}

	// fetch or create it
	obj, created, err = elasticutils.FetchOrCreateObject(proc.conn.Client, proc.conn.Config.SegmentIndex, proc.conn.Config.SegmentType, q, segFields)
	if err != nil {
		return
	}

	// and unpack the object
	err = obj.Unpack(&segobj)
	return
}
