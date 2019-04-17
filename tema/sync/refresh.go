package sync

import (
	"github.com/MathWebSearch/mwsapi/elasticutils"
)

func (proc *Process) refreshIndex() error {
	proc.print("Refreshing elasticsearch ... ")
	err := elasticutils.RefreshIndex(proc.connection.Client, proc.connection.Config.SegmentIndex, proc.connection.Config.HarvestIndex)

	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return err
}

func (proc *Process) flushIndex() error {
	proc.print("Flushing elasticsearch ... ")
	err := elasticutils.FlushIndex(proc.connection.Client, proc.connection.Config.SegmentIndex, proc.connection.Config.HarvestIndex)

	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return err
}
