package slices

import (
	"fmt"
	"reflect"
	"testing"
	
	"github.com/superwhys/venkit/lg/v2"
)

type testStruct struct {
	id   string
	name string
}

func TestDeDup(t *testing.T) {
	
	t.Run("testDedupStruct", func(t *testing.T) {
		
		ts := []testStruct{
			{id: "1", name: "one"},
			{id: "1", name: "one"},
			{id: "2", name: "two"},
			{id: "3", name: "two"},
		}
		
		deTs, err := DeDup(ts, func(idx int) string {
			return ts[idx].id
		})
		if err != nil {
			t.Errorf("DeDup() error = %v", err)
			return
		}
		lg.Info(deTs)
	})
}

func TestDupStrings(t *testing.T) {
	type args struct {
		slice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty", args{[]string{}}, []string{}},
		{"one", args{[]string{"one"}}, []string{"one"}},
		{"two-diff", args{[]string{"one", "one"}}, []string{"one"}},
		{"two-same", args{[]string{"one", "two"}}, []string{"one", "two"}},
		{"three-diff", args{[]string{"one", "one", "two"}}, []string{"one", "two"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DupStrings(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DupStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDupInt(t *testing.T) {
	type args struct {
		slice []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"empty", args{[]int{}}, []int{}},
		{"one", args{[]int{1}}, []int{1}},
		{"two-diff", args{[]int{1, 1}}, []int{1}},
		{"two-same", args{[]int{1, 2}}, []int{1, 2}},
		{"three-diff", args{[]int{1, 1, 2}}, []int{1, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DupInt(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dupint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDupSliceSmall(t *testing.T) {
	type args struct {
		slice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty", args{[]string{}}, []string{}},
		{"one", args{[]string{"one"}}, []string{"one"}},
		{"two-diff", args{[]string{"one", "one"}}, []string{"one"}},
		{"two-same", args{[]string{"one", "two"}}, []string{"one", "two"}},
		{"three-diff", args{[]string{"one", "one", "two"}}, []string{"one", "two"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DupSliceSmall(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DupSliceSmall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDupSliceLarge(t *testing.T) {
	type args struct {
		slice []string
	}
	sixty := func() []string {
		res := make([]string, 0, 60)
		for i := 0; i < 60; i++ {
			res = append(res, fmt.Sprintf("%d", i))
		}
		return res
	}()
	
	doubleSixty := func() []string {
		res := make([]string, 0, 120)
		for i := 0; i < 60; i++ {
			res = append(res, fmt.Sprintf("%d", i))
			res = append(res, fmt.Sprintf("%d", i))
		}
		return res
	}()
	
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"empty", args{[]string{}}, []string{}},
		// dedup with 60 elements which are the string of number
		{"dupLarge-1", args{sixty}, sixty},
		{"dupLarge-1", args{doubleSixty}, sixty},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dupLarge(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DupSliceLarge() = %v, want %v", got, tt.want)
			}
		})
	}
}
