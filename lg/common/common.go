package common

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func SliceClone(strSlice []string) []string {
	if strSlice == nil {
		return strSlice
	}
	return append(strSlice[:0:0], strSlice...)
}

func ParseFmtKeyValue(msg string, v ...any) (m string, keys, values []string, remains []any, err error) {
	/*
		msg= hello a=%s world %d        v = &a, 1
		hello world %d
		[true, false]
		[a]
		[%s]

		it will parse the message like a=%s out of the message and put it to the end of msg
		hello world 1 a=%s
	*/
	msgTmpl, isKv, keys, desc := ParseFmtStr(msg)
	var msgV []any
	var objV []any
	for i, kv := range isKv {
		var val any
		if i >= len(v) {
			val = "<Missing>"
		} else {
			val = v[i]
		}
		if kv {
			// a=%s
			objV = append(objV, val)
		} else {
			msgV = append(msgV, val)
		}
	}
	// fmt.Println(isKv, keys, desc, v)

	if len(isKv) < len(v) {
		remains = v[len(isKv):]
	}

	msg = fmt.Sprintf(msgTmpl, msgV...)

	if len(objV) != len(desc) {
		return "", nil, nil, nil, fmt.Errorf("invalid numbers of keys and values")
	}

	for i := range desc {
		values = append(values, fmt.Sprintf(desc[i], objV[i]))
	}

	return msg, keys, values, remains, nil
}

func ParseFmtStr(format string) (msg string, isKV []bool, keys, descs []string) {
	if json.Valid([]byte(format)) {
		return format, nil, nil, nil
	}

	var msgs []string
	for _, s := range strings.Split(format, " ") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		idx := strings.Index(s, "=%")
		if idx == -1 || strings.Contains(s[:idx], "=") {
			re, _ := regexp.Compile("%[^%]+")
			matches := re.FindAllStringIndex(s, -1)
			for i := 0; i < len(matches); i++ {
				isKV = append(isKV, false)
			}
			msgs = append(msgs, s)
			continue
		}
		keys = append(keys, s[:idx])
		descs = append(descs, s[idx+1:])
		isKV = append(isKV, true)
	}
	msg = strings.Join(msgs, " ")
	return
}
