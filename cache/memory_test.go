package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCacheStringValue(t *testing.T) {
	c := NewMemoryCache(time.Second * 10)

	c.Set("string_key", "this is a string")

	// time.Sleep(time.Second * 3)

	resp := ""
	err := c.Get("string_key", &resp)
	assert.Nil(t, err)

	assert.Equal(t, "this is a string", resp)

}
