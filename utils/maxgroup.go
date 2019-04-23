package utils

import "sync"

// MaxParallel calls worker count times with at most xthread number of element running in parallel
func MaxParallel(xthreads int, scheduler func(ch chan int64), worker func(index int64) error) (err error) {
	var ch = make(chan int64, 2*xthreads)
	var wg sync.WaitGroup
	errMutex := &sync.Mutex{}

	wg.Add(xthreads)

	for i := 0; i < xthreads; i++ {
		go func() {
			for {
				// get work or finish
				a, ok := <-ch
				if !ok {
					wg.Done()
					return
				}

				// do the work
				e := worker(a)

				// update the error if needed
				if e != nil {
					errMutex.Lock()
					if err == nil {
						err = e
					}
					defer errMutex.Unlock()
				}
			}
		}()
	}

	// Add all the work
	scheduler(ch)

	// and wait for it all to finish
	wg.Wait()
	return
}
