package wsinternal

type Subscriber interface {
	ID() string
	ExecuteCallback(msg any)
}

// TopicWorkerPool is used to divide [Subscriber]s between [TopicWorker]s.
// This prevents one [TopicWorker] of being very busy with sending messages.
type TopicWorkerPool struct {
	subscribers []Subscriber
	workers     []TopicWorker
	c           chan any
}

// NewTopicWorkerPool creates a new [TopicWorkerPool].
func NewTopicWorkerPool(maxWorkers int, channelBufferSize int) *TopicWorkerPool {
	pool := &TopicWorkerPool{
		workers:     make([]TopicWorker, maxWorkers),
		subscribers: []Subscriber{},
		c: make(
			chan any,
			channelBufferSize,
		),
	}

	pool.createWorkers(maxWorkers)

	return pool
}

// Active checks if the [TopicWorkerPool] is active
// by checking if any [TopicWorker] is active.
func (pool *TopicWorkerPool) Active() bool {
	for _, worker := range pool.workers {
		if worker.active {
			return true
		}
	}
	return false
}

// AddSubscriber adds a [Subscriber] to the [TopicWorkerPool].
func (pool *TopicWorkerPool) AddSubscriber(sub Subscriber) {
	pool.subscribers = append(pool.subscribers, sub)
}

// RemoveSubscriber removes a [Subscriber] from the [TopicWorkerPool].
func (pool *TopicWorkerPool) RemoveSubscriber(sub Subscriber) {
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
}

// Start starts a [TopicWorkerPool] if it wasn't active yet.
func (pool *TopicWorkerPool) Start() {
	if !pool.Active() {
		for i := range pool.workers {
			go pool.workers[i].HandleMessages()
		}
	}
}

// EnqueueMessage puts a message on the [TopicWorkerPool] channel.
func (pool *TopicWorkerPool) EnqueueMessage(msg any) {
	if !pool.Active() {
		return
	}

	pool.c <- msg
}

func (pool *TopicWorkerPool) createWorkers(maxWorkers int) {
	for i := 0; i < maxWorkers; i++ {
		pool.workers[i] = NewTopicWorker(i, pool)
	}
}
