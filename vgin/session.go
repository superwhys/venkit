package vgin

import (
	"encoding/gob"
	"net/url"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	redisgo "github.com/gomodule/redigo/redis"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/superwhys/venkit/v2/dialer"
)

func init() {
	gob.Register(url.Values{})
}

type RedisStoreOptions struct {
	user     string
	password string
	db       int
	keyPairs [][]byte
}

type RedisStoreOptionFunc func(o *RedisStoreOptions)

func WithUser(user string) RedisStoreOptionFunc {
	return func(o *RedisStoreOptions) {
		o.user = user
	}
}

func WithPassword(password string) RedisStoreOptionFunc {
	return func(o *RedisStoreOptions) {
		o.password = password
	}
}

func WithDb(db int) RedisStoreOptionFunc {
	return func(o *RedisStoreOptions) {
		o.db = db
	}
}

func WithKeyPairs(keyPairs ...string) RedisStoreOptionFunc {
	return func(o *RedisStoreOptions) {
		var bs [][]byte
		for _, kp := range keyPairs {
			bs = append(bs, []byte(kp))
		}
		o.keyPairs = append(o.keyPairs, bs...)
	}
}

func NewMemSessionStore(keyPairs ...string) sessions.Store {
	var bs [][]byte
	for _, kp := range keyPairs {
		bs = append(bs, []byte(kp))
	}
	return memstore.NewStore(bs...)
}

func NewCookieSessionStore(keyPairs ...string) sessions.Store {
	var bs [][]byte
	for _, kp := range keyPairs {
		bs = append(bs, []byte(kp))
	}
	return cookie.NewStore(bs...)
}

func NewRedisSessionStore(service string, opts ...RedisStoreOptionFunc) (redis.Store, error) {
	opt := &RedisStoreOptions{}

	for _, o := range opts {
		o(opt)
	}

	redisPool := dialer.DialRedisPool(service, opt.db, 100, opt.password)
	return NewRedisSessionStoreWithRedisPool(redisPool, opt.keyPairs)
}

func NewRedisSessionStoreWithRedisPool(pool *redisgo.Pool, keyPairs [][]byte) (redis.Store, error) {
	return redis.NewStoreWithPool(pool, keyPairs...)
}

func SessionDefault(ctx *Context) sessions.Session {
	return sessions.Default(ctx.Context)
}

func RegisterSessionGob(vals ...any) {
	for _, v := range vals {
		gob.Register(v)
	}
}

var (
	ErrorTokenNotFound = errors.New("Token not found!")
)

type Token interface {
	GetKey() string
	Marshal() (string, error)
	UnMarshal(val string) error
}

func SetSessionToken(c *gin.Context, t Token) error {
	session := sessions.Default(c)

	s, err := t.Marshal()
	if err != nil {
		return errors.Wrap(err, "tokenMarshal")
	}
	session.Set(t.GetKey(), s)
	if err := session.Save(); err != nil {
		return errors.Wrap(err, "saveSession")
	}

	return nil
}

func GetSessionToken(c *gin.Context, t Token) error {
	session := sessions.Default(c)

	val := session.Get(t.GetKey())
	if val == nil {
		return ErrorTokenNotFound
	}

	tokenStr, ok := val.(string)
	if !ok {
		return ErrorTokenNotFound
	}

	if err := t.UnMarshal(tokenStr); err != nil {
		return errors.Wrap(err, "decode token")
	}
	return nil
}

type StringToken string

func (st StringToken) GetKey() string {
	return string(st)
}

func (st StringToken) Marshal() (string, error) {
	return string(st), nil
}

func (st *StringToken) UnMarshal(val string) error {
	*st = StringToken(val)
	return nil
}

func NewSessionMiddleware(key string, store sessions.Store) gin.HandlerFunc {
	return sessions.Sessions(key, store)
}
