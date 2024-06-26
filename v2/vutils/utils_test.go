package vutils

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
	}{

		{"length-2", args{2}},
		{"length-4", args{4}},
		{"length-7", args{7}},
		{"length-10", args{10}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRandomString(tt.args.n)
			assert.Len(t, got, tt.args.n)
		})
	}
}

func TestMap(t *testing.T) {
	type args struct {
		arr []string
		fn  func(string) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "map-prefix",
			args: args{
				arr: []string{"hoven", "super", "value"}, fn: func(s string) string {
					return fmt.Sprintf("prefix-%v", s)
				},
			},
			want: []string{"prefix-hoven", "prefix-super", "prefix-value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.arr, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	type args struct {
		arr     []int
		initial int
		fn      func(int, int) int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test-reduce",
			args: args{
				arr:     []int{1, 2, 3},
				initial: 0,
				fn: func(u int, i int) int {
					return u + i
				},
			},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Reduce(tt.args.arr, tt.args.initial, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	type args struct {
		arr       []int
		predicate func(int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "filter-test",
			args: args{
				arr: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want: []int{2, 4, 6, 8, 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Filter(tt.args.arr, tt.args.predicate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAny(t *testing.T) {
	type args struct {
		arr       []int
		predicate func(int) bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "any-is-even-true",
			args: args{
				arr: []int{1, 2, 5, 7, 9},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want: true,
		},
		{
			name: "any-is-even-false",
			args: args{
				arr: []int{1, 5, 7, 9},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Any(tt.args.arr, tt.args.predicate); got != tt.want {
				t.Errorf("Any() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAll(t *testing.T) {
	type args struct {
		arr       []int
		predicate func(int) bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all-is-even-true",
			args: args{
				arr: []int{2, 4, 6, 8, 10},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want: true,
		},
		{
			name: "all-is-even-false",
			args: args{
				arr: []int{1, 5, 7, 9},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := All(tt.args.arr, tt.args.predicate); got != tt.want {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type args struct {
		arr       []int
		predicate func(int) bool
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name: "find-even",
			args: args{
				arr: []int{1, 3, 5, 6, 10},
				predicate: func(i int) bool {
					return i%2 == 0
				},
			},
			want:  6,
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Find(tt.args.arr, tt.args.predicate)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	type args struct {
		arr     []string
		keyFunc func(string) byte
	}
	tests := []struct {
		name string
		args args
		want map[byte][]string
	}{
		{
			name: "groupby-test",
			args: args{
				arr: []string{"super", "hoven", "hello", "summer", "hi", "yong"},
				keyFunc: func(s string) byte {
					return s[0]
				},
			},
			want: map[byte][]string{
				's': {"super", "summer"},
				'h': {"hoven", "hello", "hi"},
				'y': {"yong"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GroupBy(tt.args.arr, tt.args.keyFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZip(t *testing.T) {
	type args struct {
		arr1 []int
		arr2 []string
	}
	tests := []struct {
		name string
		args args
		want [][2]any
	}{
		{
			name: "zip-test",
			args: args{
				arr1: []int{1, 2, 3, 4, 5},
				arr2: []string{"one", "two", "three", "four", "five"},
			},
			want: [][2]any{{1, "one"}, {2, "two"}, {3, "three"}, {4, "four"}, {5, "five"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Zip(tt.args.arr1, tt.args.arr2)
			assert.Equal(t, tt.want, got)
		})
	}
}
