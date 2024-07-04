package vqueue

import (
	"fmt"
	"testing"
	"time"

	"github.com/superwhys/venkit/vredis/v2"
)

func TestMain(m *testing.M) {
	sliceQueue = NewSliceQueue[string]()

	go func() {
		for i := range 100 {
			sliceQueue.Enqueue(fmt.Sprintf("this is index %v", i))
			time.Sleep(time.Second)
		}
	}()

	redisConf := &vredis.RedisConf{
		Server: "localhost:6379",
		Db:     12,
	}
	redisQueue = NewRedisQueue[Str](redisConf, "redis-test-queue")
	go func() {
		for i := range 100 {
			if err := redisQueue.Enqueue(Str(fmt.Sprintf("this is redis index %v", i))); err != nil {
				fmt.Println("Redis Enqueue error", err)
				continue
			}
			time.Sleep(time.Second)
		}
	}()

	m.Run()
}
