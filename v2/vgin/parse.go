package vgin

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

func parseJson(_ context.Context, data []byte, params any) error {
	if err := json.Unmarshal(data, params); err != nil {
		return errors.Wrap(err, "decode json")
	}

	return nil
}

func parseMultiForm(data map[string][]string, params any) error {
	return mapFormByTag(params, data, ParamsMultiFormTag)
}

func parseHeader(needDo bool) func(c *Context, params any) error {
	return func(c *Context, params any) error {
		if !needDo || len(c.Request.Header) == 0 {
			return nil
		}

		return mapFormByTag(params, c.Request.Header, ParamsHeaderTag)
	}
}

func parsePath(needDo bool) func(c *Context, params any) error {
	return func(c *Context, params any) error {
		if !needDo || len(c.Params) == 0 {
			return nil
		}

		tmp := make(map[string][]string)
		for _, p := range c.Params {
			tmp[p.Key] = []string{p.Value}
		}
		return mapFormByTag(params, tmp, ParamsPathTag)
	}
}

func parseQuery(needDo bool) func(c *Context, params any) error {
	return func(c *Context, params any) error {
		if !needDo {
			return nil
		}

		queryMap := c.Request.URL.Query()
		return mapFormByTag(params, queryMap, ParamsQueryTag)
	}
}
