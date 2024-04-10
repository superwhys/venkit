package vflags

import "github.com/spf13/pflag"

type StringSliceGetter func() []string

func StringSlice(key string, defaultVal []string, usage string) StringSliceGetter {
	pflag.StringSlice(key, defaultVal, usage)
	v.SetDefault(key, defaultVal)
	BindPFlag(key, pflag.Lookup(key))

	return func() []string {
		return v.GetStringSlice(key)
	}
}
