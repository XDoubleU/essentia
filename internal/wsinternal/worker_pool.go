package wsinternal

import "sync"

type Subscriber interface {
	ID() string
	ExecuteCallback(msg any)
}

const subscribersUpdatedMsg = "subscribers_updated"

// WorkerPool is used to divide [Subscriber]s between [Worker]s.
// This prevents one [Worker] of being very busy.
type WorkerPool struct {
	subscribers   []Subscriber
	subscribersMu *sync.RWMutex
	workers       []Worker
}

// NewWorkerPool creates a new [WorkerPool].
func NewWorkerPool(maxWorkers int, channelBufferSize int) *WorkerPool {
	pool := &WorkerPool{
		subscribers:   []Subscriber{},
		subscribersMu: &sync.RWMutex{},
		workers:       make([]Worker, maxWorkers),
	}

	pool.createWorkers(maxWorkers, channelBufferSize)

	return pool
}

// Active checks if the [WorkerPool] is active
// by checking if any [Worker] is active.
func (pool *WorkerPool) Active() bool {
	for i := range pool.workers {
		if pool.workers[i].Active() {
			return true
		}
	}
	return false
}

// AddSubscriber adds a [Subscriber] to the [WorkerPool].
func (pool *WorkerPool) AddSubscriber(sub Subscriber) {
	pool.subscribersMu.Lock()
	defer pool.subscribersMu.Unlock()

	pool.subscribers = append(pool.subscribers, sub)
	pool.EnqueueMessage(subscribersUpdatedMsg)
}

// RemoveSubscriber removes a [Subscriber] from the [WorkerPool].
func (pool *WorkerPool) RemoveSubscriber(sub Subscriber) {
	pool.subscribersMu.Lock()
	defer pool.subscribersMu.Unlock()

	var i int
	for i = range pool.subscribers {
		if pool.subscribers[i].ID() != sub.ID() {
			continue
		}
		break
	}

	// delete subscriber
	pool.subscribers[i] = pool.subscribers[len(pool.subscribers)-1]
	pool.subscribers = pool.subscribers[:len(pool.subscribers)-1]

	pool.EnqueueMessage(subscribersUpdatedMsg)
}

// Start starts [Worker]s of a [WorkerPool] if they weren't active yet.
func (pool *WorkerPool) Start() {
	for i := range pool.workers {
		go pool.workers[i].Start()
	}
}

// EnqueueMessage puts a message on the [Worker] channels.
func (pool *WorkerPool) EnqueueMessage(msg any) {
	for i := range pool.workers {
		pool.workers[i].EnqueueMessage(msg)
	}
}

func (pool *WorkerPool) createWorkers(maxWorkers int, channelBufferSize int) {
	for i := 0; i < maxWorkers; i++ {
		pool.workers[i] = NewWorker(i, channelBufferSize, pool)
	}
}
