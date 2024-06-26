package vauth

import (
	"fmt"
	"time"
	
	"github.com/superwhys/venkit/lg/v2"
)

const (
	defaultTTL = 3600 * time.Second
)

type TokenStorage interface {
	SetWithTTL(key string, value any, ttl time.Duration) error
	Get(key string, out any) error
	Delete(key string) error
}

type Token interface {
	GetKey() string
}

type TokenManager struct {
	storage     TokenStorage
	cachePrefix string
	cacheTTL    time.Duration
}

type TokenManagerOption func(*TokenManager)

func WithCachePrefix(prefix string) TokenManagerOption {
	return func(tm *TokenManager) {
		tm.cachePrefix = prefix
	}
}

func WithCacheTTL(ttl time.Duration) TokenManagerOption {
	return func(tm *TokenManager) {
		tm.cacheTTL = ttl
	}
}

func NewTokenManager(storage TokenStorage, opts ...TokenManagerOption) *TokenManager {
	tm := &TokenManager{
		storage:  storage,
		cacheTTL: defaultTTL,
	}
	
	for _, opt := range opts {
		opt(tm)
	}
	
	return tm
}

func (tm *TokenManager) getKey(t Token) string {
	return tm.getKeyWithKeyId(t.GetKey(), t)
}

func (tm *TokenManager) getKeyWithKeyId(id string, t Token) string {
	key := tm.cachePrefix
	if key == "" {
		key = key + lg.StructName(t)
	}
	return fmt.Sprintf("%v:%v", key, id)
}

func (tm *TokenManager) Save(t Token) error {
	return tm.storage.SetWithTTL(tm.getKey(t), t, tm.cacheTTL)
}

func (tm *TokenManager) Read(key string, t Token) error {
	key = tm.getKeyWithKeyId(key, t)
	return tm.storage.Get(key, t)
}

func (tm *TokenManager) Remove(t Token) error {
	return tm.storage.Delete(tm.getKey(t))
}
