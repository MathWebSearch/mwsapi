package gogroup

import (
	"reflect"
	"testing"
	"time"
)

func TestSyncWorker(t *testing.T) {

	t.Run("linear two-group output", func(t *testing.T) {
		buffer := []int{}

		worker := newSyncWorker(1)
		worker.Work(0, func() { buffer = append(buffer, 1) })
		worker.Close(0)
		worker.Work(1, func() { buffer = append(buffer, 2) })
		worker.Close(1)
		worker.Wait()

		if !reflect.DeepEqual(buffer, []int{1, 2}) {
			t.Errorf("incorrect order")
		}
	})

	t.Run("non-linear with pause", func(t *testing.T) {
		buffer := []int{}

		worker := newSyncWorker(1)
		worker.Work(0, func() { buffer = append(buffer, 1) })
		worker.Work(1, func() { buffer = append(buffer, 3) })
		worker.Work(0, func() { buffer = append(buffer, 2) })
		time.Sleep(1 * time.Second)
		worker.Close(0)

		worker.Wait()

		if !reflect.DeepEqual(buffer, []int{1, 2, 3}) {
			t.Errorf("incorrect order")
		}
	})

	t.Run("long non-linear with pause", func(t *testing.T) {
		buffer := []int{}

		worker := newSyncWorker(1)
		worker.Work(0, func() { buffer = append(buffer, 1) })
		worker.Work(0, func() { buffer = append(buffer, 2) })
		worker.Close(0)
		worker.Work(1, func() { buffer = append(buffer, 3) })
		worker.Work(2, func() { buffer = append(buffer, 5) })
		worker.Work(3, func() { buffer = append(buffer, 7) })
		worker.Work(2, func() { buffer = append(buffer, 6) })
		worker.Close(2)
		worker.Work(1, func() { buffer = append(buffer, 4) })
		worker.Close(1)
		worker.Work(3, func() { buffer = append(buffer, 8) })
		worker.Work(4, func() { buffer = append(buffer, 9) })
		worker.Work(4, func() { buffer = append(buffer, 10) })
		worker.Close(3)
		worker.Close(4)

		worker.Wait()

		if !reflect.DeepEqual(buffer, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
			t.Errorf("SyncWorker failed to perform things in the right order")
		}
	})

}
