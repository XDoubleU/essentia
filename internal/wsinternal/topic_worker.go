package wsinternal

import (
	"math"
)

// TopicWorker is used to handle sending messages
// to a select group of [Subscriber]s of the [TopicWorkerPool].
type TopicWorker struct {
	id     int
	active bool
	pool   *TopicWorkerPool
}

// NewTopicWorker creates a new [TopicWorker].
func NewTopicWorker(id int, pool *TopicWorkerPool) TopicWorker {
	worker := TopicWorker{
		id:     id,
		active: false,
		pool:   pool,
	}
	return worker
}

// HandleMessages sends messages in the [TopicWorkerPool] channel
// to all subscribers that are handled by this [TopicWorker].
func (worker *TopicWorker) HandleMessages() {
	worker.active = true

	lowerBound, upperBound := worker.getBounds()
	for len(worker.pool.subscribers[lowerBound:upperBound]) > 0 {
		msg := <-worker.pool.c

		for _, sub := range worker.pool.subscribers[lowerBound:upperBound] {
			sub.ExecuteCallback(msg)
		}

		lowerBound, upperBound = worker.getBounds()
	}

	worker.active = false
}

func (worker *TopicWorker) getBounds() (int, int) {
	amountSubs := len(worker.pool.subscribers)
	amountWorkers := len(worker.pool.workers)

	baseBound := int(
		math.Ceil(float64(amountSubs) / float64(amountWorkers)),
	)

	lowerBound := worker.id * baseBound
	upperBound := (worker.id + 1) * baseBound

	return lowerBound, upperBound
}
