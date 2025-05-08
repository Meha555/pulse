package utils

// 由于是channel，因此内部有锁，是并发安全的
type BlockingQueue[T any] chan T

func NewBlockingQueue[T any](capacity int) *BlockingQueue[T] {
	queue := make(BlockingQueue[T], capacity)
	return &queue
}

func (t BlockingQueue[T]) Push(request T) {
	t <- request
}

func (t BlockingQueue[T]) Pop() T {
	return <-t
}

func (t BlockingQueue[T]) Len() int {
	return len(t)
}

func (t BlockingQueue[T]) Cap() int {
	return cap(t)
}

func (t BlockingQueue[T]) Close() {
	close(t)
}
