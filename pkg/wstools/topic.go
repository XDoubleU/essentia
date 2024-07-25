package wstools

import (
	"github.com/xdoubleu/essentia/internal/wsinternal"
	"nhooyr.io/websocket"
)

// Topic is used to efficiently send messages
// to [Subscriber]s in a WebSocket.
type Topic struct {
	pool               *wsinternal.TopicWorkerPool
	onSubscribeMessage any
}

// NewTopic creates a new [Topic].
func NewTopic(maxWorkers int, channelBufferSize int, onSubscribMessage any) *Topic {
	return &Topic{
		pool: wsinternal.NewTopicWorkerPool(
			maxWorkers,
			channelBufferSize,
		),
		onSubscribeMessage: onSubscribMessage,
	}
}

// Subscribe subscribes a [Subscriber] to this [Topic].
// If configured a message will be sent on subscribing.
// If no message handling go routine was
// running this will be started now.
func (t *Topic) Subscribe(conn *websocket.Conn) {
	sub := NewSubscriber(t, conn)
	t.pool.AddSubscriber(sub)

	if t.onSubscribeMessage != nil {
		sub.ExecuteCallback(t.onSubscribeMessage)
	}

	t.pool.Start()
}

// UnSubscribe unsubscribes a [Subscriber] from this [Topic].
func (t *Topic) UnSubscribe(sub Subscriber) {
	t.pool.RemoveSubscriber(sub)
}

// EnqueueMessage enqueues a message if there are subscribers on this [Topic].
func (t *Topic) EnqueueMessage(msg any) {
	t.pool.EnqueueMessage(msg)
}
