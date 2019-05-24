package utils

import (
	"reflect"
	"testing"
)

func TestFilterInt64(t *testing.T) {
	type args struct {
		slice  []int64
		filter func(int64) bool
	}
	tests := []struct {
		name    string
		args    args
		wantRes []int64
	}{
		{"keep all elements", args{slice: []int64{1, 2, 3, 4, 5}, filter: func(_ int64) bool { return true }}, []int64{1, 2, 3, 4, 5}},
		{"keep no elements", args{slice: []int64{1, 2, 3, 4, 5}, filter: func(_ int64) bool { return false }}, []int64{}},
		{"keep even element", args{slice: []int64{1, 2, 3, 4, 5}, filter: func(i int64) bool { return i%2 == 0 }}, []int64{2, 4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := FilterInt64(tt.args.slice, tt.args.filter); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("FilterInt64() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestContainsInt64(t *testing.T) {
	type args struct {
		haystack []int64
		needle   int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Does not contain element", args{[]int64{1, 2, 3, 4, 5}, 6}, false},
		{"Does contain element", args{[]int64{1, 2, 3, 4, 5}, 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsInt64(tt.args.haystack, tt.args.needle); got != tt.want {
				t.Errorf("ContainsInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
