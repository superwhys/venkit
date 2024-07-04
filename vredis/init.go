package vredis

import (
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/dialer"
	"github.com/superwhys/venkit/v2/vflags"
)

var redisConfFlag = vflags.Struct("redisConf", (*RedisConf)(nil), "Redis config")

type RedisConf struct {
	Server   string `desc:"redis server name (default localhost:6379)"`
	Password string `desc:"redis server password"`
	Db       int    `desc:"redis db (default 0)"`
	MaxIdle  int    `desc:"redis maxIdle (default 100)"`
}

func (conf *RedisConf) DialRedisPool() *redis.Pool {
	return dialer.DialRedisPool(
		conf.Server,
		conf.Db,
		conf.MaxIdle,
		conf.Password,
	)
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

var RedisConn = func() *RedisClient {
	lg.Fatal("vredis.RedisConn is not initialize, you can set `VENKIT_AUTO_REDIS=1` environment variable to enable automatic redis connection.")
	return nil
}

func init() {
	if os.Getenv(autoRedisKey) != "1" {
		return
	}

	conf := &RedisConf{}
	lg.PanicError(redisConfFlag(conf))

	var pwd []string
	if conf.Password != "" {
		pwd = append(pwd, conf.Password)
	}

	lg.Debugf("auto connect to redis with config: %v", lg.Jsonify(conf))
	RedisConn = func() *RedisClient {
		return NewRedisClient(conf.DialRedisPool())
	}
}
