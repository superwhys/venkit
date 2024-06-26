package vhttp

import (
	"encoding/json"
	
	"github.com/superwhys/venkit/lg/v2"
)

type JsonBody map[string]any

func NewJsonBody() JsonBody {
	return make(JsonBody)
}

func (j JsonBody) Add(key string, value any) JsonBody {
	_, exists := j[key]
	if !exists {
		j[key] = value
	}
	return j
}

func (j JsonBody) Encode() []byte {
	b, err := json.Marshal(j)
	if err != nil {
		lg.Errorf("json body encode error: %v", err)
		return nil
	}
	
	return b
}
