package vflags

import "github.com/spf13/pflag"

type IntGetter func() int

func Int(key string, defaultVal int, usage string) IntGetter {
	pflag.Int(key, defaultVal, usage)
	v.SetDefault(key, defaultVal)
	BindPFlag(key, pflag.Lookup(key))

	return func() int {
		return v.GetInt(key)
	}
}

func IntRequired(key, usage string) IntGetter {
	requiredFlags = append(requiredFlags, key)
	return Int(key, 0, usage)
}
