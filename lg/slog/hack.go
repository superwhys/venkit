package slog

import (
	"log/slog"
	_ "unsafe"
)

//go:linkname argsToAttrSlice log/slog.argsToAttrSlice
func argsToAttrSlice(args []any) []slog.Attr
