package vflags

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/superwhys/venkit/lg"
)

var (
	ErrNotStruct = errors.New("not struct")
)

type HasDefault interface {
	SetDefault()
}

type HasValidator interface {
	Validate() error
}

func Struct(key string, defaultVal any, usage string) func(out any) error {
	err := setPFlagRecursively(key, defaultVal)
	if err != nil {
		if !errors.Is(err, ErrNotStruct) {
			lg.Errorf("Bind struct flags error: %v", err)
			lg.PanicError(err)
		}
		lg.Warnf("it won't display `%v` desciption with not struct default val", key)
	}
	v.SetDefault(key, defaultVal)
	return func(out any) error {
		if err := v.UnmarshalKey(key, out); err != nil {
			return err
		}

		if d, ok := out.(HasDefault); ok {
			d.SetDefault()
		}
		if v, ok := out.(HasValidator); ok {
			return v.Validate()
		}

		return nil
	}
}

func setPFlag(key string, ptr interface{}) {
	v.BindPFlag(key, pflag.Lookup(key))
	nestedKey[key] = ptr
}

func setPFlagRecursively(prefix string, i interface{}) error {
	vf := reflect.ValueOf(i)
	if vf.Kind() == reflect.Ptr {
		vf = vf.Elem()
	}
	if vf.Kind() != reflect.Struct {
		return ErrNotStruct
	}
	for i := 0; i < vf.NumField(); i++ {
		field := vf.Type().Field(i)
		name := field.Name
		for _, tag := range []string{"vflags", "json"} {
			if content := field.Tag.Get(tag); content != "" {
				name = strings.SplitN(content, ",", 2)[0]
				break
			}
		}
		usage := field.Tag.Get("usage")
		name = prefix + "." + name

		switch vf.Field(i).Kind() {
		case reflect.String:
			setPFlag(name, pflag.String(name, vf.Field(i).String(), usage))
		case reflect.Bool:
			setPFlag(name, pflag.Bool(name, vf.Field(i).Bool(), usage))
		case reflect.Int, reflect.Int64:
			if field.Type.String() == "time.Duration" {
				setPFlag(name, pflag.Duration(name, time.Duration(vf.Field(i).Int()), usage))
			} else {
				setPFlag(name, pflag.Int(name, int(vf.Field(i).Int()), usage))
			}
		case reflect.Float64:
			setPFlag(name, pflag.Float64(name, vf.Field(i).Float(), usage))
		case reflect.Slice:
			switch field.Type.String() {
			case "[]int":
				setPFlag(name, pflag.IntSlice(name, vf.Field(i).Interface().([]int), usage))
			case "[]string":
				setPFlag(name, pflag.StringSlice(name, vf.Field(i).Interface().([]string), usage))
			case "[]float64":
				setPFlag(name, pflag.Float64Slice(name, vf.Field(i).Interface().([]float64), usage))
			case "[]bool":
				setPFlag(name, pflag.BoolSlice(name, vf.Field(i).Interface().([]bool), usage))
			case "[]time.Duration":
				setPFlag(name, pflag.DurationSlice(name, vf.Field(i).Interface().([]time.Duration), usage))
			default:
				return fmt.Errorf("unsupport type of field %s %s", field.Name, field.Type.String())
			}
		case reflect.Struct, reflect.Ptr:
			if err := setPFlagRecursively(name, vf.Field(i).Interface()); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupport kind of field %s %s", field.Name, vf.Field(i).Kind())
		}
	}

	return nil
}

func injectNestedKey() {
	for key, valuePtr := range nestedKey {
		flag := pflag.Lookup(key)
		if flag != nil && flag.Changed {
			switch valuePtr.(type) {
			case *int:
				v.Set(key, *valuePtr.(*int))
			case *bool:
				v.Set(key, *valuePtr.(*bool))
			case *float64:
				v.Set(key, *valuePtr.(*float64))
			case *time.Duration:
				v.Set(key, *valuePtr.(*time.Duration))
			case *string:
				v.Set(key, *valuePtr.(*string))
			case *[]bool:
				v.Set(key, *valuePtr.(*[]bool))
			case *[]string:
				v.Set(key, *valuePtr.(*[]string))
			case *[]int:
				v.Set(key, *valuePtr.(*[]int))
			case *[]float64:
				v.Set(key, *valuePtr.(*[]float64))
			case *[]time.Duration:
				v.Set(key, *valuePtr.(*[]time.Duration))
			default:
				lg.Fatal("Unsupport flag value type", flag.Value.Type())
			}
		}
	}
}
