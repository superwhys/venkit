package vqueue

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/superwhys/venkit/vredis/v2"
)

var _ Queue[*noItem] = (*RedisQueue[*noItem])(nil)

type Item interface {
	Key() string
}

type noItem struct{}

func (i *noItem) Key() string {
	return ""
}

type RedisQueue[T Item] struct {
	*vredis.RedisClient
	queue string
}

func NewRedisQueue[T Item](conf *vredis.RedisConf, queue string) *RedisQueue[T] {
	pool := conf.DialRedisPool()

	return &RedisQueue[T]{
		RedisClient: vredis.NewRedisClient(pool),
		queue:       queue,
	}
}

func (q *RedisQueue[T]) Enqueue(value T) error {
	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = q.Do("RPUSH", q.queue, serializedValue)
	return err
}

func (q *RedisQueue[T]) parseQueueData(bulk []any) (T, error) {
	var ret T
	var zero T
	if err := json.Unmarshal(bulk[1].([]byte), &ret); err != nil {
		return zero, errors.Wrap(err, "redisDequeueDecode")
	}

	return ret, nil
}

func (q *RedisQueue[T]) Dequeue() (T, error) {
	var zero T

	bulks, err := redis.Values(q.Do("BLPOP", q.queue, 5))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return zero, QueueEmptyError
		}
		return zero, err
	}

	return q.parseQueueData(bulks)
}

func (q *RedisQueue[T]) size() (int, error) {
	length, err := redis.Int(q.Do("LLEN", q.queue))
	if err != nil {
		return 0, err
	}

	return length, nil
}

func (q *RedisQueue[T]) IsEmpty() (bool, error) {
	length, err := q.size()
	if err != nil {
		return false, err
	}

	return length == 0, nil
}

func (q *RedisQueue[T]) Size() (int, error) {
	return q.size()
}
