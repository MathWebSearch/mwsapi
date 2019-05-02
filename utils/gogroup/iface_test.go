package gogroup

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewWorkGroupEngine(t *testing.T) {
	type args struct {
		xthreads  int
		needsSync bool
	}
	tests := []struct {
		name       string
		args       args
		wantEngine string
	}{
		{"parallel engine (1)", args{10, false}, "parallel"},
		{"parallel engine (2)", args{10, true}, "parallel"},

		{"async engine (1)", args{0, false}, "async"},
		{"async engine (2)", args{0, true}, "async"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewWorkGroup(tt.args.xthreads, tt.args.needsSync)
			defer group.Close()
			if gotEngine := group.Engine(); !reflect.DeepEqual(gotEngine, tt.wantEngine) {
				t.Errorf("NewWorkGroup().Engine() = %v, want %v", gotEngine, tt.wantEngine)
			}
		})
	}
}

func TestNewWorkGroupJobs(t *testing.T) {
	type args struct {
		xthreads  int
		needsSync bool
		jobs      int
		runner    func(int, func(i int)) error
	}
	tests := []struct {
		name string
		args args

		wantErr         bool
		wantMaxParallel int
		wantMissedJobs  []int
		wantBuffers     map[int][]int
	}{
		{"100 jobs with 10 run in parallel", args{
			xthreads: 10,

			jobs: 100,
			runner: func(i int, _ func(i int)) error {
				time.Sleep(500 * time.Millisecond)
				return nil
			},
		}, false, 10, nil, make(map[int][]int)},

		{"100 jobs, with unlimited number of parallel jobs", args{
			xthreads: 0,

			jobs: 100,
			runner: func(i int, _ func(i int)) error {
				time.Sleep(500 * time.Millisecond)
				return nil
			},
		}, false, 100, nil, make(map[int][]int)},

		{"2 jobs (in sequence) returning errors", args{
			xthreads: 1,

			jobs: 2,
			runner: func(i int, _ func(i int)) error {
				time.Sleep(500 * time.Millisecond)
				return errors.New("Something went wrong")
			},
		}, true, 1, nil, make(map[int][]int)},

		{"2 jobs (in parallel) returning errors", args{
			xthreads: 0,

			jobs: 2,
			runner: func(i int, _ func(i int)) error {
				time.Sleep(500 * time.Millisecond)
				return errors.New("Something went wrong")
			},
		}, true, 2, nil, make(map[int][]int)},

		{"2 sequential jobs write the right buffer", args{
			xthreads:  1,
			needsSync: true,

			jobs: 2,
			runner: func(i int, call func(i int)) error {
				time.Sleep(500 * time.Millisecond)
				call(2*i + 1)
				time.Sleep(500 * time.Millisecond)
				call(2*i + 2)
				return nil
			},
		}, false, 1, nil, map[int][]int{0: []int{1, 2}, 1: []int{3, 4}}},

		{"2 parallel jobs write the right buffer", args{
			xthreads:  0,
			needsSync: true,

			jobs: 2,
			runner: func(i int, call func(i int)) error {
				time.Sleep(500 * time.Millisecond)
				call(2*i + 1)
				time.Sleep(500 * time.Millisecond)
				call(2*i + 2)
				return nil
			},
		}, false, 2, nil, map[int][]int{0: []int{1, 2}, 1: []int{3, 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewWorkGroup(tt.args.xthreads, tt.args.needsSync)
			res := runJobs(group, tt.args.jobs, tt.args.runner)

			if (res.err != nil) != tt.wantErr {
				t.Errorf("NewWorkGroup() error = %v, wantErr %v", res.err, tt.wantErr)
			}

			if res.maxParallel != tt.wantMaxParallel {
				t.Errorf("NewWorkGroup() maxParallel = %v, want %v", res.maxParallel, tt.wantMaxParallel)
			}

			if !reflect.DeepEqual(res.missedJobs, tt.wantMissedJobs) {
				t.Errorf("NewWorkGroup() missedJobs = %v, want %v", res.missedJobs, tt.wantMissedJobs)
			}

			if !reflect.DeepEqual(res.buffers, tt.wantBuffers) {
				t.Errorf("NewWorkGroup() buffers = %v, want %v", res.buffers, tt.wantBuffers)
			}
		})
	}
}

type groupTestResult struct {
	maxParallel int
	err         error

	missedJobs []int
	buffers    map[int][]int
}

// testWorkGroup runs tests in a work group and keeps track of results
func runJobs(group WorkGroup, jobs int, runner func(int, func(i int)) error) (res *groupTestResult) {
	res = &groupTestResult{
		buffers: make(map[int][]int),
	}

	ranJobs := make([]bool, jobs)

	// current number of jobs running
	numJobs := 0
	jobmutex := &sync.Mutex{}

	// quenue all the jobs
	for i := 0; i < jobs; i++ {
		func(i int) {
			ranJobs[i] = false

			job := GroupJob(func(sync func(func())) error {
				// we are starting this job, so check if we need to add it to the parallel job count
				jobmutex.Lock()
				numJobs++
				if numJobs > res.maxParallel {
					res.maxParallel = numJobs
				}
				jobmutex.Unlock()

				// decrease the number of jobs when done
				defer func() {
					jobmutex.Lock()
					numJobs--
					jobmutex.Unlock()
				}()

				// store that we ran the job
				defer func() {
					ranJobs[i] = true
				}()

				// append the job to the buffer
				return runner(i, func(j int) {
					sync(func() {
						res.buffers[i] = append(res.buffers[i], j)
					})
				})
			})

			group.Add(&job)
		}(i)
	}

	// wait for the group
	res.err = group.Wait()

	// and find the missing jobs
	for i, ran := range ranJobs {
		if !ran {
			res.missedJobs = append(res.missedJobs, i)
		}
	}

	return
}
