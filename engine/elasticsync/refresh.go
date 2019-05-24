package elasticsync

import (
	"github.com/MathWebSearch/mwsapi/utils/elasticutils"
	"github.com/pkg/errors"
)

func (proc *Process) refreshIndex() (err error) {
	proc.Print(nil, "Refreshing elasticsearch ... ")
	err = elasticutils.RefreshIndex(proc.conn.Client, proc.conn.Config.SegmentIndex, proc.conn.Config.HarvestIndex)
	err = errors.Wrap(err, "elasticutils.RefreshIndex failed")

	if err == nil {
		proc.PrintlnOK(nil, "OK")
	} else {
		proc.PrintlnERROR(nil, "ERROR")
	}

	return
}

func (proc *Process) flushIndex() (err error) {
	proc.Print(nil, "Flushing elasticsearch ... ")
	err = elasticutils.FlushIndex(proc.conn.Client, proc.conn.Config.SegmentIndex, proc.conn.Config.HarvestIndex)
	err = errors.Wrap(err, "elasticutils.FlushIndex failed")

	if err == nil {
		proc.PrintlnOK(nil, "OK")
	} else {
		proc.PrintlnERROR(nil, "ERROR")
	}

	return
}
