package wsinternal

import (
	"context"
	"math"
	"sync"
)

// Worker is used to handle work of the [WorkerPool].
type Worker struct {
	id         int
	active     bool
	activeMu   *sync.RWMutex
	lowerBound int
	upperBound int
	c          chan any
	pool       *WorkerPool
}

// NewWorker creates a new [Worker].
func NewWorker(id int, channelBufferSize int, pool *WorkerPool) Worker {
	worker := Worker{
		id:         id,
		active:     false,
		activeMu:   &sync.RWMutex{},
		lowerBound: -1,
		upperBound: -1,
		c: make(
			chan any,
			channelBufferSize,
		),
		pool: pool,
	}
	return worker
}

func (worker *Worker) Active() bool {
	worker.activeMu.RLock()
	defer worker.activeMu.RUnlock()

	return worker.active
}

func (worker *Worker) EnqueueEvent(event any) {
	if !worker.Active() {
		return
	}

	worker.c <- event
}

// Start makes [Worker] start doing work.
func (worker *Worker) Start(_ context.Context) error {
	// already active
	if worker.Active() {
		return nil
	}

	// if lock not free is either being checked to start or is being started
	if !worker.activeMu.TryLock() {
		return nil
	}

	worker.active = true
	worker.activeMu.Unlock()

	worker.calculateBounds()

	for {
		// no subscribers so stop worker
		// if lock is free no one is currently checking current state
		if worker.upperBound == 0 && worker.activeMu.TryLock() {
			worker.active = false
			worker.activeMu.Unlock()
			break
		}

		event := <-worker.c

		// stop has been called from pool
		if event == stopEvent {
			worker.activeMu.Lock()
			worker.active = false
			worker.activeMu.Unlock()
			break
		}

		// subscribers have been updated, have to update bounds
		if event == subscribersUpdatedEvent {
			worker.calculateBounds()
			continue
		}

		worker.pool.subscribersMu.RLock()

		// no work so check again later
		if worker.upperBound > len(worker.pool.subscribers) {
			worker.pool.subscribersMu.RUnlock()
			continue
		}

		for _, sub := range worker.pool.subscribers[worker.lowerBound:worker.upperBound] {
			sub.OnEventCallback(event)
		}

		worker.pool.subscribersMu.RUnlock()
	}

	return nil
}

func (worker *Worker) calculateBounds() {
	worker.pool.subscribersMu.RLock()
	defer worker.pool.subscribersMu.RUnlock()

	amountSubs := len(worker.pool.subscribers)
	amountWorkers := len(worker.pool.workers)

	subsPerWorker := int(
		math.Ceil(float64(amountSubs) / float64(amountWorkers)),
	)

	worker.lowerBound = worker.id * subsPerWorker
	worker.upperBound = (worker.id + 1) * subsPerWorker
}
