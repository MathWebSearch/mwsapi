package elasticsync

import (
	"sync/atomic"

	"github.com/MathWebSearch/mwsapi/utils/gogroup"

	"github.com/MathWebSearch/mwsapi/utils"
)

// upsertSegments updates or inserts new segements
func (proc *Process) upsertSegments() (err error) {
	proc.Println(nil, "Upserting harvest segments ...")

	group := gogroup.NewWorkGroup(proc.conn.Config.PoolSize, true)

	err = utils.IterateFiles(proc.harvestFolder, ".json", func(path string) error {
		job := gogroup.GroupJob(func(sync func(func())) error {
			proc.Printf(sync, "=> %s\n", path)
			return proc.upsertSegment(sync, path)
		})
		group.Add(&job)
		return nil
	})

	err = group.UWait(err)

	if err == nil {
		proc.PrintlnOK(nil, "OK")
	} else {
		proc.PrintlnERROR(nil, "ERROR")
	}

	return
}

// upsertSegment upserts a single segment
func (proc *Process) upsertSegment(sync func(func()), segment string) (err error) {
	// compute the hash
	proc.Print(sync, "  computing hash ... ")
	hash, err := utils.HashFile(segment)

	if err != nil {
		proc.PrintlnERROR(sync, "ERROR")
		return err
	}
	proc.Printf(sync, "%s\n", hash)

	// check the index
	proc.Print(sync, "  checking segment index ... ")

	seg, obj, created, err := proc.checkSegmentIndex(segment)
	if err != nil {
		proc.PrintlnERROR(sync, "ERROR")
		return err
	}

	if created {
		proc.PrintlnSKIP(sync, "NOT FOUND")
	} else {
		proc.PrintlnOK(sync, "FOUND")
	}

	proc.Print(sync, "  Comparing segment hash ... ")

	hashMatch := seg.Hash == hash
	if hashMatch {
		proc.PrintlnOK(sync, "MATCHES")
	} else {
		proc.PrintlnSKIP(sync, "DIFFERS")
	}

	// if the hash matches, we don't need to update
	if !hashMatch || proc.force {
		if hashMatch && proc.force {
			proc.Println(sync, "  Hash matches, but --force was given. Forcing update. ")
		}

		if created {
			atomic.AddInt64(&proc.stats.NewSegments, 1)
		} else {
			atomic.AddInt64(&proc.stats.UpdatedSegments, 1)
		}

		proc.Print(sync, "  Clearing harvests belonging to segment ... ")
		err = proc.clearSegmentHarvests(seg)
		if err != nil {
			proc.PrintlnERROR(sync, "ERROR")
			return err
		}
		proc.PrintlnOK(sync, "OK")

		// we need to clear out the old segments from the db, and put the new ones in
		proc.Print(sync, "  Loading harvests from segment into index ... ")
		err = proc.insertSegmentHarvests(segment)
		if err != nil {
			proc.PrintlnERROR(sync, "ERROR")
			return err
		}
		proc.PrintlnOK(sync, "OK")
	} else {
		atomic.AddInt64(&proc.stats.UnchangedSegments, 1)
	}

	proc.Print(sync, "  Storing segment state ... ")
	seg.Touched = true
	seg.Hash = hash

	// repack and save it
	err = obj.Pack(seg)
	if err == nil {
		err = obj.Save()
	}

	if err == nil {
		proc.PrintlnOK(sync, "OK")
	} else {
		proc.PrintlnERROR(sync, "ERROR")
	}

	return err
}
