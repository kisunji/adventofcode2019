package main

import (
	"testing"
)

func Test_calculateFuel(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "12 yields 2",
			args: args{12},
			want: 2,
		},
		{
			name: "14 yields 2",
			args: args{14},
			want: 2,
		},
		{
			name: "1969 yields 654",
			args: args{1969},
			want: 966,
		},
		{
			name: "100756 yields 2",
			args: args{100756},
			want: 50346,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateFuel(tt.args.i); got != tt.want {
				t.Errorf("calculateFuel() = %v, want %v", got, tt.want)
			}
		})
	}
}
