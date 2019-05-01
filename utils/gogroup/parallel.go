package gogroup

import (
	"sync"
)

// parallelGroup represents a group that can run an unlimited number of groups in parallel
type parallelGroup struct {
	group     *sync.WaitGroup // internal wait group for quenued jobs
	workGroup *sync.WaitGroup // group for adding
	xthreads  int             // maximum number of parallel threads

	worker    *syncWorker // sync worker
	needsSync bool        // do we need the worker?

	ch chan *GroupJob // the channel to jobs

	indexMutex *sync.Mutex // for protecting the index
	index      int64       // counter for the index

	errMutex *sync.Mutex // for protecting the error
	err      error       // the error itself
}

// newParallelGroup makes a new parallel group
func newParallelGroup(xthreads int, needsSync bool) (group *parallelGroup) {
	// create a new group
	group = &parallelGroup{
		group:     &sync.WaitGroup{},
		workGroup: &sync.WaitGroup{},
		xthreads:  xthreads,

		needsSync: needsSync,

		ch: make(chan *GroupJob, xthreads*2),

		indexMutex: &sync.Mutex{},

		errMutex: &sync.Mutex{},
	}

	// create a worker if we need it
	if needsSync {
		group.worker = newSyncWorker()
	}

	// and start the group
	group.start()

	return
}

func (group *parallelGroup) start() {
	group.group.Add(group.xthreads)

	// add all the threads
	for i := 0; i < group.xthreads; i++ {
		go func() {
			defer group.group.Done()
			for {
				// get work or finish
				job, ok := <-group.ch
				if !ok {
					return
				}

				// get the index of this job
				group.indexMutex.Lock()
				index := group.index
				group.index++
				group.indexMutex.Unlock()

				// do the work
				e := (*job)(func(sync func()) {
					if group.needsSync {
						group.worker.Work(index, sync)
					}
				})

				if group.needsSync {
					group.worker.Close(index)
				}

				// update the error if needed
				if e != nil {
					group.errMutex.Lock()
					if group.err == nil {
						group.err = e
					}
					group.errMutex.Unlock()
				}
			}
		}()
	}

}

// Engine returns the name of the engine for debugging
func (group *parallelGroup) Engine() string {
	return "parallel"
}

// Add schedules an extra job to run
func (group *parallelGroup) Add(job *GroupJob) {
	group.workGroup.Add(1)
	go func() {
		defer group.workGroup.Done()
		group.ch <- job
	}()
}

// Wait waits for this group to finish
func (group *parallelGroup) Wait() (err error) {
	group.workGroup.Wait()
	close(group.ch)
	group.group.Wait()
	if group.needsSync {
		group.worker.Wait()
	}
	return group.err
}

// UWait is the same as wait, but only updates error if it isn't nil
func (group *parallelGroup) UWait(err error) error {
	e := group.Wait()
	if err == nil {
		err = e
	}
	return err
}
