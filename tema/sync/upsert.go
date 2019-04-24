package sync

import (
	"sync/atomic"

	"github.com/MathWebSearch/mwsapi/utils/gogroup"

	"github.com/MathWebSearch/mwsapi/utils"
)

// upsertSegments updates or inserts new segements
func (proc *Process) upsertSegments() (err error) {
	proc.println(nil, "Upserting harvest segments ...")

	group := gogroup.NewWorkGroup(proc.connection.Config.PoolSize, true)

	err = utils.IterateFiles(proc.harvestFolder, ".json", func(path string) error {
		job := gogroup.GroupJob(func(sync func(func())) error {
			proc.printf(sync, "=> %s\n", path)
			return proc.upsertSegment(sync, path)
		})
		group.Add(&job)
		return nil
	})

	err = group.UWait(err)

	if err == nil {
		proc.printlnOK(nil, "OK")
	} else {
		proc.printlnERROR(nil, "ERROR")
	}

	return
}

// upsertSegment upserts a single segment
func (proc *Process) upsertSegment(sync func(func()), segment string) (err error) {
	// compute the hash
	proc.print(sync, "  computing hash ... ")
	hash, err := utils.HashFile(segment)

	if err != nil {
		proc.printlnERROR(sync, "ERROR")
		return err
	}
	proc.printf(sync, "%s\n", hash)

	// check the index
	proc.print(sync, "  checking segment index ... ")

	seg, obj, created, err := proc.checkSegmentIndex(segment)
	if err != nil {
		proc.printlnERROR(sync, "ERROR")
		return err
	}

	if created {
		proc.printlnSKIP(sync, "NOT FOUND")
	} else {
		proc.printlnOK(sync, "FOUND")
	}

	proc.print(sync, "  Comparing segment hash ... ")

	hashMatch := seg.Hash == hash
	if hashMatch {
		proc.printlnOK(sync, "MATCHES")
	} else {
		proc.printlnSKIP(sync, "DIFFERS")
	}

	// if the hash matches, we don't need to update
	if !hashMatch || proc.force {
		if hashMatch && proc.force {
			proc.println(sync, "  Hash matches, but --force was given. Forcing update. ")
		}

		if created {
			atomic.AddInt64(&proc.stats.NewSegments, 1)
		} else {
			atomic.AddInt64(&proc.stats.UpdatedSegments, 1)
		}

		proc.print(sync, "  Clearing harvests belonging to segment ... ")
		err = proc.clearSegmentHarvests(seg)
		if err != nil {
			proc.printlnERROR(sync, "ERROR")
			return err
		}
		proc.printlnOK(sync, "OK")

		// we need to clear out the old segments from the db, and put the new ones in
		proc.print(sync, "  Loading harvests from segment into index ... ")
		err = proc.insertSegmentHarvests(segment)
		if err != nil {
			proc.printlnERROR(sync, "ERROR")
			return err
		}
		proc.printlnOK(sync, "OK")
	} else {
		atomic.AddInt64(&proc.stats.UnchangedSegments, 1)
	}

	proc.print(sync, "  Storing segment state ... ")
	seg.Touched = true
	seg.Hash = hash

	// repack and save it
	err = obj.Pack(seg)
	if err == nil {
		err = obj.Save()
	}

	if err == nil {
		proc.printlnOK(sync, "OK")
	} else {
		proc.printlnERROR(sync, "ERROR")
	}

	return err
}
