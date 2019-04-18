package utils

import "sync"

// AsyncGroup represents a group that can manage errors
type AsyncGroup struct {
	group  *sync.WaitGroup
	worker *SyncWorker

	indexMutex *sync.Mutex
	index      int

	errMutex *sync.Mutex
	err      error
}

// NewAsyncGroup makes a new AsyncGroup
func NewAsyncGroup() *AsyncGroup {
	return &AsyncGroup{
		group:  &sync.WaitGroup{},
		worker: NewSyncWorker(1),

		indexMutex: &sync.Mutex{},
		errMutex:   &sync.Mutex{},
	}
}

// Add adds a job to this group
func (group *AsyncGroup) Add(job func(func(func())) error) *AsyncGroup {
	group.indexMutex.Lock()
	defer group.indexMutex.Unlock()

	group.group.Add(1)

	go func(i int) {
		defer group.group.Done()

		err := job(func(sync func()) {
			group.worker.Work(i, sync)
		})

		if err != nil {
			group.errMutex.Lock()
			if group.err == nil {
				group.err = err
			}
			group.errMutex.Unlock()
		}

		group.worker.Close(i)
	}(group.index)

	group.index++
	return group
}

// Wait waits for this group to finish
func (group *AsyncGroup) Wait() (err error) {
	group.group.Wait()
	group.worker.Wait()
	return group.err
}

// UWait is the same as wait, but only updates error if it isn't nil
func (group *AsyncGroup) UWait(err error) error {
	e := group.Wait()
	if err == nil {
		err = e
	}
	return err
}
