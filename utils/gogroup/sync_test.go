package gogroup

import (
	"reflect"
	"testing"
	"time"
)

func TestSyncWorker(t *testing.T) {
	type args struct {
		runner func(*syncWorker, func(i int) func())
	}
	tests := []struct {
		name      string
		args      args
		wantOrder []int
	}{
		{"linear two-group-output", args{func(worker *syncWorker, call func(i int) func()) {
			worker.SendWork(0, call(1))
			worker.SendClose(0)
			worker.SendWork(1, call(2))
			worker.SendClose(1)
		}}, []int{1, 2}},

		{"non-linear with pause", args{func(worker *syncWorker, call func(i int) func()) {
			worker.SendWork(0, call(1))
			worker.SendWork(1, call(3))
			worker.SendWork(0, call(2))
			time.Sleep(1 * time.Second)
			worker.SendClose(0)
		}}, []int{1, 2, 3}},

		{"long non-linear without pause", args{func(worker *syncWorker, call func(i int) func()) {
			worker.SendWork(0, call(1))
			worker.SendWork(0, call(2))
			worker.SendClose(0)
			worker.SendWork(1, call(3))
			worker.SendWork(2, call(5))
			worker.SendWork(3, call(7))
			worker.SendWork(2, call(6))
			worker.SendClose(2)
			worker.SendWork(1, call(4))
			worker.SendClose(1)
			worker.SendWork(3, call(8))
			worker.SendWork(4, call(9))
			worker.SendWork(4, call(10))
			worker.SendClose(3)
			worker.SendClose(4)
		}}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOrder := []int{}

			worker := newSyncWorker()
			tt.args.runner(worker, func(i int) func() {
				return func() {
					gotOrder = append(gotOrder, i)
				}
			})
			worker.Wait()

			if !reflect.DeepEqual(gotOrder, tt.wantOrder) {
				t.Errorf("newSyncWorker() call order = %#v, expected %#v", gotOrder, tt.wantOrder)
			}
		})
	}
}
