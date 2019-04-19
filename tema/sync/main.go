package sync

import (
	"fmt"
	"time"

	"github.com/MathWebSearch/mwsapi/tema"
)

// Process represent args to the syncronisation process
type Process struct {
	connection *tema.Connection

	force         bool
	harvestFolder string

	quiet bool
	stats *Stats
}

// Stats represents the stats of a process
type Stats struct {
	UnchangedSegments int64 // Segments which were not changed
	UpdatedSegments   int64 // Segments which were updated
	NewSegments       int64 // Segments which were newly added
	RemovedSegments   int64 // Segements which were removed

	Duration time.Duration // how long it took
}

func (s *Stats) String() string {
	return fmt.Sprintf("took %s: new %d, updated %d, unchanged %d, removed %d", s.Duration, s.NewSegments, s.UpdatedSegments, s.UnchangedSegments, s.RemovedSegments)
}

// NewProcess creates a new Process
func NewProcess(connection *tema.Connection, folder string, quiet bool, force bool) *Process {
	return &Process{
		connection:    connection,
		harvestFolder: folder,

		force: force,
		quiet: quiet,
	}
}

// Run is the main sync entry point
func (proc *Process) Run() (stats *Stats, err error) {
	// reset stats
	proc.stats = &Stats{}

	// keep track of how long the process takes
	start := time.Now()
	defer func() {
		stats.Duration = time.Since(start)
	}()

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
