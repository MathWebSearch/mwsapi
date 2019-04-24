package gogroup

// WorkGroup represents a group that performs parallel work
type WorkGroup interface {
	Add(job *GroupJob) WorkGroup // Adds a job to this group
	Wait() error                 // Waits for this group to finish, returning any errors
	UWait(err error) error       // Same as wait, but only returns an error iff the given one is not nil
}

// NewWorkGroup creates a new work group
// xthreads is the maximal number of parallel jobs <= 0 for unlimited.
// needsSync should be true if jobs also need to perform syncronous work, such as logging
func NewWorkGroup(xthreads int, needsSync bool) (group WorkGroup) {
	if xthreads <= 0 {
		group = newAsyncGroup(needsSync)
	} else {
		group = newParallelGroup(xthreads, needsSync)
	}

	return
}

// GroupJob represents the job for a group
type GroupJob func(func(func())) error
