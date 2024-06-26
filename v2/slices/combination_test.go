package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombinationsSlice(t *testing.T) {
	type args struct {
		list []string
		n    int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{"comb-2", args{
			list: []string{"a", "b", "c"},
			n:    2,
		}, [][]string{
			{"a", "b"},
			{"a", "c"},
			{"b", "c"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CombinationsSlice(tt.args.list, tt.args.n)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRandomSelect(t *testing.T) {
	type args struct {
		source []string
		count  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"random-select-1", args{
			source: []string{"a", "b", "c", "d"},
			count:  1,
		}, false},
		{"random-select-2", args{
			source: []string{"a", "b", "c", "d"},
			count:  2,
		}, false},
		{"random-select-5", args{
			source: []string{"a", "b", "c", "d"},
			count:  5,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomSelect(tt.args.source, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("RandomSelect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Len(t, got, tt.args.count)
		})
	}
}
