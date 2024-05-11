package vgin

import (
	"bytes"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	bodyBufferKey = "venkit.vgin.bodyBuffer"
	bodyDataKey   = "venkit.vgin.body"
)

var bufferPool = &sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func BodyBufferMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		c.Set(bodyBufferKey, buf)

		c.Next()

		o, exists := c.Get(bodyBufferKey)
		if exists {
			buf = o.(*bytes.Buffer)
			c.Set(bodyBufferKey, nil)
			bufferPool.Put(buf)
			c.Request.Body = nil
		}
	}
}

func BodyRawData(c *gin.Context) ([]byte, error) {
	if b, exists := c.Get(bodyDataKey); exists {
		return b.([]byte), nil
	}

	var buf *bytes.Buffer

	if o, exists := c.Get(bodyBufferKey); exists {
		buf = o.(*bytes.Buffer)
	} else {
		buf = new(bytes.Buffer)
	}

	body, err := c.GetRawData()
	if err != nil {
		return nil, err
	}

	if _, err = buf.Write(body); err != nil {
		return nil, err
	}

	c.Set(bodyDataKey, body)

	c.Request.Body = io.NopCloser(buf)
	return body, nil
}
