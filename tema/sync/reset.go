package sync

import (
	"github.com/MathWebSearch/mwsapi/elasticutils"
	"gopkg.in/olivere/elastic.v6"
)

// resetSegmentIndex resets the segment index to prepare for updates
func (proc *Process) resetSegmentIndex() (err error) {
	proc.print(nil, "Resetting Segment Index ... ")

	// reset the touched part to false
	err = elasticutils.UpdateAll(proc.connection.Client, proc.connection.Config.SegmentIndex, proc.connection.Config.SegmentType, elastic.NewScript("ctx._source.touched = false"))
	if err == nil {
		proc.printlnOK(nil, "OK")
	} else {
		proc.printlnERROR(nil, "ERROR")
	}

	return
}
