package gogroup

import (
	"sync"

	"github.com/MathWebSearch/mwsapi/utils"
)

// parallelGroup represents a group that can run an unlimited number of groups in parallel
type parallelGroup struct {
	group    *sync.WaitGroup // internal wait group for quenued jobs
	xthreads int             // maximum number of parallel threads

	worker    *utils.SyncWorker // sync worker
	needsSync bool              // do we need the worker?

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
		group:    &sync.WaitGroup{},
		xthreads: xthreads,

		needsSync: needsSync,

		ch: make(chan *GroupJob, xthreads*2),

		indexMutex: &sync.Mutex{},

		errMutex: &sync.Mutex{},
	}

	// create a worker if we need it
	if needsSync {
		group.worker = utils.NewSyncWorker(1)
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
			for {
				// get work or finish
				job, ok := <-group.ch
				if !ok {
					group.group.Done()
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
					defer group.errMutex.Unlock()
				}
			}
		}()
	}

}

// Add schedules an extra job to run
func (group *parallelGroup) Add(job *GroupJob) WorkGroup {
	group.ch <- job
	return group
}

// Wait waits for this group to finish
func (group *parallelGroup) Wait() (err error) {
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
