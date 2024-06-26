package dialer

import (
	"time"
	
	"github.com/gomodule/redigo/redis"
	"github.com/superwhys/venkit/lg/v2"
	"github.com/superwhys/venkit/v2/discover"
)

func DialRedisPool(addr string, db int, maxIdle int, password ...string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: 300 * time.Second,
		Dial:        consulRedisDial(addr, db, password...),
	}
}

func consulRedisDial(addr string, db int, password ...string) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		var serviceAddr string
		serviceAddr = discover.GetServiceFinder().GetAddress(addr)
		if serviceAddr == "" {
			serviceAddr = addr
		}
		lg.Debugf("Discover redis addr: %v", serviceAddr)
		
		options := []redis.DialOption{
			redis.DialDatabase(db),
			redis.DialConnectTimeout(5 * time.Second),
		}
		
		if len(password) > 0 && password[0] != "" {
			options = append(options, redis.DialPassword(password[0]))
		}
		
		return redis.Dial("tcp", serviceAddr, options...)
	}
}
