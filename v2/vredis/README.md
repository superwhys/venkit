# vredis

This is a simple encapsulation of `redis` and provides a variety of common methods like `Get`, `Set`, `Delete`,`Lock`, `Unlock`

## Example
You can create a new `RedisClient` by `NewRedisClient()`

```Go
func NewRedisClient(pool *redis.Pool) *RedisClient {
	return &RedisClient{
		pool: pool,
	}
}
```

You can also set the `VENKIT_AUTO_REDIS` environment variable to enable automatic connection to `RedisClient`

but it is important to note that you need to specify a `configuration file` that includes the configuration required for redis connection
```Go
type RedisConf struct {
	Server   string `desc:"redis server name (default localhost:6379)"`
	Password string `desc:"redis server password"`
	Db       int    `desc:"redis db (default 0)"`
	MaxIdle  int    `desc:"redis maxIdle (default 100)"`
}
```

the configuration may be like
```yaml
...
redisConf:
  server: localhost:6379
  password: redispwd
  db: 10
  maxIdle: 100
...
```

