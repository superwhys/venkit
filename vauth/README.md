# Vauth

## Example
```go
package main

import (
	"github.com/google/uuid"
	"github.com/superwhys/venkit/cache/v2"
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/dialer"
)

type TestToken struct {
	Uid   string
	Value string
}

func (t *TestToken) GetKey() string {
	return t.Uid
}

func main() {
	redisCache := cache.NewRedisCache(dialer.DialRedisPool("localhost:6379", 14, 100))
	tm := NewTokenManager(redisCache)

	token := &TestToken{Uid: uuid.NewString(), Value: "testtokenvalue"}

	if err := tm.Save(token); err != nil {
		t.Error(err)
		return
	}

	newToken := &TestToken{}
	if err := tm.Read(token.Uid, newToken); err != nil {
		t.Error(err)
		return
	}

	if newToken.Value != token.Value {
		lg.Errorf("not equal, get: %v, want: %v", lg.Jsonify(newToken), lg.Jsonify(token))
		return
	}
}

```
