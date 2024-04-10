package authutils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/superwhys/venkit/cache"
	"github.com/superwhys/venkit/dialer"
	"github.com/superwhys/venkit/lg"
)

type TestToken struct {
	Uid   string
	Value string
}

func (t *TestToken) GetKey() string {
	return t.Uid
}

func (t *TestToken) SetKey(val string) {
	t.Uid = val
}

func TestTokenManager(t *testing.T) {
	redisCache := cache.NewRedisCache(dialer.DialRedisPool("localhost:6379", 14, 100))
	tm := NewTokenManager(redisCache)

	token := &TestToken{Uid: uuid.NewString(), Value: "testtokenvalue"}

	if err := tm.Save(token); err != nil {
		t.Error(err)
		return
	}

	newToken := &TestToken{Uid: token.Uid}
	if err := tm.Read(newToken); err != nil {
		t.Error(err)
		return
	}

	if newToken.Value != token.Value {
		t.Errorf("not equal, get: %v, want: %v", lg.Jsonify(newToken), lg.Jsonify(token))
		return
	}
}
