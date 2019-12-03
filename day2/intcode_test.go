package main

import (
	"reflect"
	"testing"
)

func Test_computeInput(t *testing.T) {
	type args struct {
		arr []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "first small program",
			args: args{[]int{1, 0, 0, 0, 99}},
			want: []int{2, 0, 0, 0, 99},
		},
		{
			name: "second small program",
			args: args{[]int{2, 3, 0, 3, 99}},
			want: []int{2, 3, 0, 6, 99},
		},
		{
			name: "third small program",
			args: args{[]int{2, 4, 4, 5, 99, 0}},
			want: []int{2, 4, 4, 5, 99, 9801},
		},
		{
			name: "fourth small program",
			args: args{[]int{1, 1, 1, 4, 99, 5, 6, 0, 99}},
			want: []int{30, 1, 1, 4, 2, 5, 6, 0, 99},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := computeInput(tt.args.arr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
