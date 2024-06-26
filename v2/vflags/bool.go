package vflags

import "github.com/spf13/pflag"

type BoolGetter func() bool

func Bool(key string, defaultVal bool, usage string) BoolGetter {
	pflag.Bool(key, defaultVal, usage)
	v.SetDefault(key, defaultVal)
	BindPFlag(key, pflag.Lookup(key))

	return func() bool {
		return v.GetBool(key)
	}
}

func BoolRequired(key, usage string) BoolGetter {
	requiredFlags = append(requiredFlags, key)
	return Bool(key, false, usage)
}
