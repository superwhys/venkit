package vqueue

var _ Queue[any] = (*SliceQueue[any])(nil)

type SliceQueue[T any] struct {
	items []T
}

func NewSliceQueue[T any]() *SliceQueue[T] {
	return &SliceQueue[T]{items: make([]T, 0)}
}

func (q *SliceQueue[T]) Enqueue(value T) error {
	q.items = append(q.items, value)
	return nil
}

func (q *SliceQueue[T]) Dequeue() (T, error) {
	if len(q.items) == 0 {
		var zero T
		return zero, QueueEmptyError
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

func (q *SliceQueue[T]) IsEmpty() (bool, error) {
	return len(q.items) == 0, nil
}

func (q *SliceQueue[T]) Size() (int, error) {
	return len(q.items), nil
}
