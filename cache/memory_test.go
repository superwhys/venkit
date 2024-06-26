package cache

import (
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/venkit/lg/v2"
)

type People interface {
	SayHello() string
}

type Peter struct {
	name     string
	timeFunc func() time.Time
}

func (p *Peter) SayHello() string {
	return p.name + "" + "hello" + "" + p.timeFunc().String()
}

func TestMemoryCacheStringValue(t *testing.T) {
	c := NewMemoryCache(time.Second * 10)
	
	c.Set("string_key", "this is a string")
	
	// time.Sleep(time.Second * 3)
	
	resp := ""
	err := c.Get("string_key", &resp)
	assert.Nil(t, err)
	
	assert.Equal(t, "this is a string", resp)
}

func getPeople(now time.Time) People {
	return &Peter{
		name: "everyone",
		timeFunc: func() time.Time {
			return now
		},
	}
}

func TestMemoryCacheInterfaceValue(t *testing.T) {
	now := time.Now()
	c := NewMemoryCache(time.Second * 10)
	
	c.Set("interface_key", getPeople(now))
	
	var resp People
	err := c.Get("interface_key", &resp)
	assert.Nil(t, err)
	
	assert.Equal(t, getPeople(now).SayHello(), resp.SayHello())
	lg.Info(resp.SayHello())
}
