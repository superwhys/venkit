package vredis

import (
	"os"

	"github.com/superwhys/venkit/dialer"
	"github.com/superwhys/venkit/lg"
	"github.com/superwhys/venkit/vflags"
)

type RedisConf struct {
	Server   string `desc:"redis server name (default localhost:6379)"`
	Password string `desc:"redis server password"`
	Db       int    `desc:"redis db (default 0)"`
	MaxIdle  int    `desc:"redis maxIdle (default 100)"`
}

func (rc *RedisConf) SetDefault() {
	if rc.Server == "" && rc.MaxIdle == 0 {
		rc.Server = "localhost:6379"
		rc.MaxIdle = 100
	}
}

var (
	autoRedisKey = "VENKIT_AUTO_REDIS"
)

var RedisConn func() *RedisClient

func init() {
	if os.Getenv(autoRedisKey) != "1" {
		return
	}

	redisConfFlag := vflags.Struct("redisConf", &RedisConf{}, "Redis config")

	RedisConn = func() *RedisClient {
		conf := &RedisConf{}
		lg.PanicError(redisConfFlag(conf))

		var pwd []string
		if conf.Password != "" {
			pwd = append(pwd, conf.Password)
		}

		lg.Debugf("auto connect to redis with config: %v", lg.Jsonify(conf))
		return NewRedisClient(dialer.DialRedisPool(
			conf.Server,
			conf.Db,
			conf.MaxIdle,
			pwd...,
		))
	}
}
