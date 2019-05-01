package gogroup

import (
	"sync"
)

// asyncGroup represents a group that can run an unlimited number of groups in parallel
type asyncGroup struct {
	group *sync.WaitGroup // internal wait group for quenued jobs

	worker    *syncWorker // sync worker
	needsSync bool        // do we need the worker?

	indexMutex *sync.Mutex // for protecting the index
	index      int64       // counter for the index

	errMutex *sync.Mutex // for protecting the error
	err      error       // the error itself
}

// NewAsyncGroup makes a new AsyncGroup
func newAsyncGroup(needsSync bool) (group *asyncGroup) {

	// create a new group
	group = &asyncGroup{
		group:     &sync.WaitGroup{},
		needsSync: needsSync,

		indexMutex: &sync.Mutex{},
		errMutex:   &sync.Mutex{},
	}

	// create a worker if we need it
	if needsSync {
		group.worker = newSyncWorker()
	}

	return
}

// Engine returns the name of the engine for debugging
func (group *asyncGroup) Engine() string {
	return "async"
}

// Add adds a job to this group
func (group *asyncGroup) Add(job *GroupJob) {
	group.indexMutex.Lock()
	defer group.indexMutex.Unlock()

	group.group.Add(1)

	go func(i int64) {
		defer group.group.Done()

		err := (*job)(func(sync func()) {
			if group.needsSync {
				group.worker.Work(i, sync)
			}
		})

		if err != nil {
			group.errMutex.Lock()
			if group.err == nil {
				group.err = err
			}
			group.errMutex.Unlock()
		}

		if group.needsSync {
			group.worker.Close(int64(i))
		}
	}(group.index)

	group.index++
}

// Wait waits for this group to finish
func (group *asyncGroup) Wait() (err error) {
	group.group.Wait()
	if group.needsSync {
		group.worker.Wait()
	}
	return group.err
}

// UWait is the same as wait, but only updates error if it isn't nil
func (group *asyncGroup) UWait(err error) error {
	e := group.Wait()
	if err == nil {
		err = e
	}
	return err
}
