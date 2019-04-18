package sync

import (
	"github.com/MathWebSearch/mwsapi/elasticutils"
)

func (proc *Process) refreshIndex() error {
	proc.print(nil, "Refreshing elasticsearch ... ")
	err := elasticutils.RefreshIndex(proc.connection.Client, proc.connection.Config.SegmentIndex, proc.connection.Config.HarvestIndex)

	if err == nil {
		proc.printlnOK(nil, "OK")
	} else {
		proc.printlnERROR(nil, "ERROR")
	}

	return err
}

func (proc *Process) flushIndex() error {
	proc.print(nil, "Flushing elasticsearch ... ")
	err := elasticutils.FlushIndex(proc.connection.Client, proc.connection.Config.SegmentIndex, proc.connection.Config.HarvestIndex)

	if err == nil {
		proc.printlnOK(nil, "OK")
	} else {
		proc.printlnERROR(nil, "ERROR")
	}

	return err
}
