//go:build !windows
// +build !windows

package vflags

import "syscall"

func kill() {
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
