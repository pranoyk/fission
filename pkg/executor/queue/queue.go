package queue

import (
	"errors"
)

type Queue struct {
	identity string
	queue    chan interface{}
	lock     chan struct{}
}

func NewQueue(capacity int, identity string) *Queue {
	queue := &Queue{}
	queue.initialize(capacity, identity)

	return queue
}

func (q *Queue) initialize(capacity int, identity string) {
	q.queue = make(chan interface{}, capacity)
	q.identity = identity
	q.lock = make(chan struct{}, 1)
}

func (q *Queue) IsValidIdentity(identity string) bool {
	return q.identity == identity
}

func (q *Queue) Lock() {
	// non-blocking fill the channel
	select {
	case q.lock <- struct{}{}:
	default:
	}
}

func (q *Queue) Unlock() {
	// non-blocking flush the channel
	select {
	case <-q.lock:
	default:
	}
}

func (q *Queue) GetCap() int {
	q.Lock()
	defer q.Unlock()

	return cap(q.queue)
}

func (q *Queue) IsLocked() bool {
	return len(q.lock) >= 1
}

func (q *Queue) Enqueue(value interface{}) error {
	if q.IsLocked() {
		return errors.New("the queue is locked")
	}
	select {
	case q.queue <- value:
	default:
		return errors.New("Queue is at full capacity")
	}
	return nil
}

func (q *Queue) Dequeue() (interface{}, error) {
	if q.IsLocked() {
		return nil, errors.New("the queue is locked")
	}

	select {
	case value, ok := <-q.queue:
		if ok {
			return value, nil
		}
		return nil, errors.New("internal channel is closed")
	default:
		return nil, errors.New("empty queue")
	}
}
