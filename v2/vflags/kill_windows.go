//go:build windows
// +build windows

package vflags

import (
	"os"
)

func kill() {
	os.Exit(1)
}
