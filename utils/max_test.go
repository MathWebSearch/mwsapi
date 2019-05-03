package utils

import (
	"testing"
)

func TestMaxInt(t *testing.T) {
	type args struct {
		values []int
	}
	tests := []struct {
		name    string
		args    args
		wantMax int
	}{
		{"maximum of empty list", args{values: []int{}}, 0},

		{"maximum of one-element positive list", args{values: []int{1}}, 1},
		{"maximum of one-element negative list", args{values: []int{-1}}, -1},

		{"maximum of positive numbers", args{values: []int{1, 2, 3, 4, 5}}, 5},
		{"maximum of negative numbers", args{values: []int{-1, -2, -3, -4, -5}}, -1},

		{"maximum of mixed set", args{values: []int{1, -3, 47, 333, 12, -78}}, 333},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMax := MaxInt(tt.args.values...); gotMax != tt.wantMax {
				t.Errorf("MaxInt() = %v, want %v", gotMax, tt.wantMax)
			}
		})
	}
}

func TestMaxInt64(t *testing.T) {
	type args struct {
		values []int64
	}
	tests := []struct {
		name    string
		args    args
		wantMax int64
	}{
		{"maximum of empty list", args{values: []int64{}}, 0},

		{"maximum of one-element positive list", args{values: []int64{1}}, 1},
		{"maximum of one-element negative list", args{values: []int64{-1}}, -1},

		{"maximum of positive numbers", args{values: []int64{1, 2, 3, 4, 5}}, 5},
		{"maximum of negative numbers", args{values: []int64{-1, -2, -3, -4, -5}}, -1},

		{"maximum of mixed set", args{values: []int64{1, -3, 47, 333, 12, -78}}, 333}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMax := MaxInt64(tt.args.values...); gotMax != tt.wantMax {
				t.Errorf("MaxInt64() = %v, want %v", gotMax, tt.wantMax)
			}
		})
	}
}
