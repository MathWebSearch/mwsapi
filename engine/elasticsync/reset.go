package elasticsync

import (
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
	"gopkg.in/olivere/elastic.v6"
)

// resetSegmentIndex resets the segment index to prepare for updates
func (proc *Process) resetSegmentIndex() (err error) {
	proc.Print(nil, "Resetting Segment Index ... ")

	// reset the touched part to false
	err = elasticutils.UpdateAll(proc.conn.Client, proc.conn.Config.SegmentIndex, proc.conn.Config.SegmentType, elastic.NewScript("ctx._source.touched = false"))
	if err == nil {
		proc.PrintlnOK(nil, "OK")
	} else {
		proc.PrintlnERROR(nil, "ERROR")
	}

	return
}
