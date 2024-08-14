package ws

import (
	"context"

	"github.com/xdoubleu/essentia/internal/wsinternal"
	"nhooyr.io/websocket"
)

// OnSubscribeCallback is called to fetch data that
// should be returned when a new subscriber is added to a topic.
type OnSubscribeCallback = func(ctx context.Context, topic *Topic) (any, error)

// Topic is used to efficiently send messages
// to [Subscriber]s in a WebSocket.
type Topic struct {
	Name                string
	pool                *wsinternal.WorkerPool
	onSubscribeCallback OnSubscribeCallback
}

// NewTopic creates a new [Topic].
func NewTopic(
	name string,
	maxWorkers int,
	channelBufferSize int,
	onSubscribeCallback OnSubscribeCallback,
) *Topic {
	return &Topic{
		Name: name,
		pool: wsinternal.NewWorkerPool(
			maxWorkers,
			channelBufferSize,
		),
		onSubscribeCallback: onSubscribeCallback,
	}
}

// Subscribe subscribes a [Subscriber] to this [Topic].
// If configured a message will be sent on subscribing.
// If no message handling go routine was
// running this will be started now.
func (t *Topic) Subscribe(conn *websocket.Conn) error {
	sub := NewSubscriber(t, conn)
	t.pool.AddSubscriber(sub)

	if t.onSubscribeCallback != nil {
		event, err := t.onSubscribeCallback(context.Background(), t)
		if err != nil {
			return err
		}

		sub.OnEventCallback(event)
	}

	t.pool.Start()

	return nil
}

// UnSubscribe unsubscribes a [Subscriber] from this [Topic].
func (t *Topic) UnSubscribe(sub Subscriber) {
	t.pool.RemoveSubscriber(sub)
}

// EnqueueEvent enqueues an event if there are subscribers on this [Topic].
func (t *Topic) EnqueueEvent(event any) {
	t.pool.EnqueueEvent(event)
}

// StopPool stops the used [wsinternal.WorkerPool].
func (t *Topic) StopPool() {
	t.pool.Stop()
}
