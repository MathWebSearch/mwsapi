package sync

import (
	"fmt"

	"github.com/MathWebSearch/mwsapi/tema"
)

// upsertSegments updates or inserts new segements
func (proc *Process) upsertSegments() (err error) {
	proc.println("Upserting harvest segments ...")

	err = iterateFiles(proc.harvestFolder, ".json", func(path string) error {
		proc.printf("=> %s\n", path)
		return proc.upsertSegment(path)
	})

	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return
}

// upsertSegment upserts a single segment
func (proc *Process) upsertSegment(segment string) (err error) {
	// compute the hash
	proc.print("  computing hash ... ")
	hash, err := tema.HashSegment(segment)

	if err != nil {
		proc.printlnERROR("ERROR")
		return err
	}
	proc.printf("%s\n", hash)

	// check the index
	proc.print("  checking segment index ... ")

	seg, obj, created, err := proc.checkSegmentIndex(segment)
	if err != nil {
		proc.printlnERROR("ERROR")
		return err
	}

	if created {
		proc.printlnOK("NOT FOUND")
	} else {
		proc.printlnOK("FOUND")
	}

	proc.print("  Comparing segment hash ... ")

	// if the hash matches, we don't need to update
	if seg.Hash != hash {
		proc.printlnOK("DIFFERS")

		if created {
			proc.stats.NewSegments++
		} else {
			proc.stats.UpdatedSegments++
		}

		proc.print("  Clearing harvests belonging to segment ... ")
		err = proc.clearSegmentHarvests(seg)
		if err != nil {
			proc.printlnERROR("ERROR")
			return err
		}
		proc.printlnOK("OK")

		// we need to clear out the old segments from the db, and put the new ones in
		fmt.Print("  Loading harvests from segment into index ... ")
		err = proc.insertSegmentHarvests(segment)
		if err != nil {
			proc.printlnERROR("ERROR")
			return err
		}
		proc.printlnOK("OK")
	} else {
		proc.printlnOK("MATCHES")
		proc.stats.UnchangedSegments++
	}

	proc.print("  Storing segment state ... ")
	seg.Touched = true
	seg.Hash = hash

	// repack and save it
	err = obj.Pack(seg)
	if err == nil {
		err = obj.Save()
	}

	if err == nil {
		proc.printlnOK("OK")
	} else {
		proc.printlnERROR("ERROR")
	}

	return err
}
