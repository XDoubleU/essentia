package ws

import (
	"github.com/xdoubleu/essentia/internal/wsinternal"
	"nhooyr.io/websocket"
)

// OnSubscribeCallback is called to fetch data that
// should be returned when a new subscriber is added to a topic.
type OnSubscribeCallback = func(topic *Topic) any

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
func (t *Topic) Subscribe(conn *websocket.Conn) {
	sub := NewSubscriber(t, conn)
	t.pool.AddSubscriber(sub)

	if t.onSubscribeCallback != nil {
		event := t.onSubscribeCallback(t)
		sub.OnEventCallback(event)
	}

	t.pool.Start()
}

// UnSubscribe unsubscribes a [Subscriber] from this [Topic].
func (t *Topic) UnSubscribe(sub Subscriber) {
	t.pool.RemoveSubscriber(sub)
}

// EnqueueEvent enqueues an event if there are subscribers on this [Topic].
func (t *Topic) EnqueueEvent(msg any) {
	t.pool.EnqueueEvent(msg)
}

// StopPool stops the used [wsinternal.WorkerPool].
func (t *Topic) StopPool() {
	t.pool.Stop()
}
