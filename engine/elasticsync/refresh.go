package elasticsync

import (
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
)

func (proc *Process) refreshIndex() error {
	proc.Print(nil, "Refreshing elasticsearch ... ")
	err := elasticutils.RefreshIndex(proc.conn.Client, proc.conn.Config.SegmentIndex, proc.conn.Config.HarvestIndex)

	if err == nil {
		proc.PrintlnOK(nil, "OK")
	} else {
		proc.PrintlnERROR(nil, "ERROR")
	}

	return err
}

func (proc *Process) flushIndex() error {
	proc.Print(nil, "Flushing elasticsearch ... ")
	err := elasticutils.FlushIndex(proc.conn.Client, proc.conn.Config.SegmentIndex, proc.conn.Config.HarvestIndex)

	if err == nil {
		proc.PrintlnOK(nil, "OK")
	} else {
		proc.PrintlnERROR(nil, "ERROR")
	}

	return err
}
