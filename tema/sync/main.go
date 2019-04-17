package sync

import (
	"github.com/MathWebSearch/mwsapi/tema"
)

// Process represent args to the syncronisation process
type Process struct {
	connection *tema.Connection

	harvestFolder string
}

// NewProcess creates a new Process
func NewProcess(connection *tema.Connection, Folder string) *Process {
	return &Process{
		connection:    connection,
		harvestFolder: Folder,
	}
}

// Run is the main sync entry point
func (proc *Process) Run() {
	// Create the index and mapping
	err := proc.createIndex()
	if err != nil {
		panic(err)
	}

	// Reset the segment index
	err = proc.resetSegmentIndex()
	if err != nil {
		panic(err)
	}

	// refresh all the indexes
	err = proc.refreshIndex()
	if err != nil {
		panic(err)
	}

	// upsert segments
	err = proc.upsertSegments()
	if err != nil {
		panic(err)
	}

	// refresh all the indexes
	err = proc.refreshIndex()
	if err != nil {
		panic(err)
	}

	// clear old segements
	err = proc.clearSegments()
	if err != nil {
		panic(err)
	}

	// flush all the indexes
	err = proc.flushIndex()
	if err != nil {
		panic(err)
	}

	// and be done
	return
}
