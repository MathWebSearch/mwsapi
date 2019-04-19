package sync

import "github.com/MathWebSearch/mwsapi/elasticutils"

// createIndex creates an index to prepare for segmented updates
func (proc *Process) createIndex() (err error) {
	proc.printf(nil, "Creating Harvest Index %s ... ", proc.connection.Config.HarvestIndex)
	created, err := elasticutils.CreateIndex(proc.connection.Client, proc.connection.Config.HarvestIndex, proc.connection.Config.HarvestMapping())
	if err != nil {
		proc.printlnERROR(nil, "ERROR")
		return
	}
	if created {
		proc.printlnOK(nil, "OK")
	} else {
		proc.printlnSKIP(nil, "SKIP")
	}

	proc.printf(nil, "Creating Segment Index %s ... ", proc.connection.Config.SegmentIndex)
	created, err = elasticutils.CreateIndex(proc.connection.Client, proc.connection.Config.SegmentIndex, proc.connection.Config.SegmentMapping())
	if err != nil {
		proc.printlnERROR(nil, "ERROR")
		return
	}
	if created {
		proc.printlnOK(nil, "OK")
	} else {
		proc.printlnSKIP(nil, "SKIP")
	}

	return
}
