package vqueue

import "errors"

var (
	QueueEmptyError = errors.New("queue is empty")
)

type Queue[T any] interface {
	Enqueue(value T) error
	Dequeue() (T, error)
	IsEmpty() (bool, error)
	Size() (int, error)
}
