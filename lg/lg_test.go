package lg

import (
	"fmt"
	"testing"
)

func TestInfo(t *testing.T) {
	tests := []struct {
		name string
		args string
	}{
		{"test1", "hello info"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args)
		})
	}
}

func TestDebug(t *testing.T) {
	tests := []struct {
		name string
		args string
	}{
		{"test1", "hello debug"},
	}
	EnableDebug()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debug(tt.args)
		})
	}
}

func TestTimeFuncDefer(t *testing.T) {
	fn := func() {
		fmt.Println("before defer")
		defer TimeDurationDefer()()

		fmt.Println("after defer")
	}

	fn()
}
