package dialer

import (
	"github.com/superwhys/venkit/dialer"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/snail"
	"github.com/superwhys/venkit/vflags"
	"github.com/superwhys/venkit/vredis"
)

type RedisConf struct {
	Server   string `desc:"redis server name (default localhost:6379)"`
	Password string `desc:"redis server password"`
	Db       int    `desc:"redis db (default 0)"`
	MaxIdle  int    `desc:"redis maxIdle (default 100)"`
}

func (rc *RedisConf) SetDefault() {
	rc.Server = "localhost:6379"
	rc.Db = 0
	rc.MaxIdle = 100
}

var (
	redisConfFlag = vflags.Struct("redisConf", &RedisConf{}, "Redis config")
)

var Client *vredis.RedisClient

func init() {
	snail.RegisterObject("redisClient", func() error {
		conf := &RedisConf{}
		lg.PanicError(redisConfFlag(conf))

		var pwd []string
		if conf.Password != "" {
			pwd = append(pwd, conf.Password)
		}
		Client = vredis.NewRedisClient(dialer.DialRedisPool(
			conf.Server,
			conf.Db,
			conf.MaxIdle,
			pwd...,
		))
		return nil
	})
}
