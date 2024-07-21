package wstools

import (
	"context"

	"nhooyr.io/websocket"
)

// Topic is used to efficiently send messages
// to [Subscriber]s in a WebSocket.
type Topic struct {
	pool               *TopicWorkerPool
	onSubscribeMessage any
}

// NewTopic creates a new [Topic].
func NewTopic(maxWorkers int, channelBufferSize int, onSubscribMessage any) *Topic {
	return &Topic{
		pool:               NewTopicWorkerPool(maxWorkers, channelBufferSize),
		onSubscribeMessage: onSubscribMessage,
	}
}

// Subscribe subscribes a [Subscriber] to this [Topic].
// If configured a message will be sent on subscribing.
// If no message handling go routine was
// running this will be started now.
func (t *Topic) Subscribe(ctx context.Context, conn *websocket.Conn) {
	sub := t.pool.AddSubscriber(ctx, t, conn)

	if t.onSubscribeMessage != nil {
		sub.SendMessage(t.onSubscribeMessage)
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
