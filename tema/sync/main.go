package sync

import (
	"fmt"

	"github.com/MathWebSearch/mwsapi/tema"
)

// Process represent args to the syncronisation process
type Process struct {
	connection *tema.Connection

	quiet bool
	stats *Stats

	harvestFolder string
}

// Stats represents the stats of a process
type Stats struct {
	UnchangedSegments int64 // Segments which were not changed
	UpdatedSegments   int64 // Segments which were updated
	NewSegments       int64 // Segments which were newly added
	RemovedSegments   int64 // Segements which were removed
}

func (s *Stats) String() string {
	return fmt.Sprintf("new %d, updated %d, unchanged %d, removed %d", s.NewSegments, s.UpdatedSegments, s.UnchangedSegments, s.RemovedSegments)
}

// NewProcess creates a new Process
func NewProcess(connection *tema.Connection, folder string, quiet bool) *Process {
	return &Process{
		connection:    connection,
		harvestFolder: folder,
		quiet:         quiet,
	}
}

// Run is the main sync entry point
func (proc *Process) Run() (stats *Stats, err error) {
	// reset stats
	proc.stats = &Stats{}

	// Create the index and mapping
	err = proc.createIndex()
	if err != nil {
		return
	}

	// Reset the segment index
	err = proc.resetSegmentIndex()
	if err != nil {
		return
	}

	// refresh all the indexes
	err = proc.refreshIndex()
	if err != nil {
		return
	}

	// upsert segments
	err = proc.upsertSegments()
	if err != nil {
		return
	}

	// refresh all the indexes
	err = proc.refreshIndex()
	if err != nil {
		return
	}

	// clear old segements
	err = proc.clearSegments()
	if err != nil {
		return
	}

	// flush all the indexes
	err = proc.flushIndex()
	if err != nil {
		return
	}

	// and return with stats
	stats = proc.stats
	return
}
