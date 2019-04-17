package sync

import (
	"github.com/MathWebSearch/mwsapi/tema"

	"github.com/MathWebSearch/mwsapi/elasticutils"
	"gopkg.in/olivere/elastic.v6"
)

// checkSegmentIndex checks the segment index for a given segment
func (proc *Process) checkSegmentIndex(segment string) (segobj *tema.Segment, obj *elasticutils.Object, created bool, err error) {
	// the query
	q := elastic.NewBoolQuery()
	q = q.Must(elastic.NewTermQuery("segment", segment))

	// the fields
	segFields := tema.Segment{
		ID:      segment,
		Hash:    "",
		Touched: true,
	}

	// serialize the new fields
	fields := make(map[string]interface{})
	err = elasticutils.Repack(segFields, &fields)
	if err != nil {
		return
	}

	// fetch or create it
	obj, created, err = elasticutils.FetchOrCreateObject(proc.connection.Client, proc.connection.Config.SegmentIndex, proc.connection.Config.SegmentType, q, fields)
	if err != nil {
		return
	}

	// and unpack the object
	err = obj.Unpack(&segobj)
	return
}
