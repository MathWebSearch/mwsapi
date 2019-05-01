package elasticsync

import "github.com/MathWebSearch/mwsapi/utils/elasticutils"

// createIndex creates an index to prepare for segmented updates
func (proc *Process) createIndex() (err error) {
	proc.Printf(nil, "Creating Harvest Index %s ... ", proc.conn.Config.HarvestIndex)
	created, err := elasticutils.CreateIndex(proc.conn.Client, proc.conn.Config.HarvestIndex, harvestMapping(proc.conn.Config))
	if err != nil {
		proc.PrintlnERROR(nil, "ERROR")
		return
	}
	if created {
		proc.PrintlnOK(nil, "OK")
	} else {
		proc.PrintlnSKIP(nil, "SKIP")
	}

	proc.Printf(nil, "Creating Segment Index %s ... ", proc.conn.Config.SegmentIndex)
	created, err = elasticutils.CreateIndex(proc.conn.Client, proc.conn.Config.SegmentIndex, segmentMapping(proc.conn.Config))
	if err != nil {
		proc.PrintlnERROR(nil, "ERROR")
		return
	}
	if created {
		proc.PrintlnOK(nil, "OK")
	} else {
		proc.PrintlnSKIP(nil, "SKIP")
	}

	return
}
