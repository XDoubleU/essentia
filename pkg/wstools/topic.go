package wstools

import (
	"context"
	"sync"

	"nhooyr.io/websocket"
)

// Topic is used to efficiently send messages
// to [Subscriber]s in a WebSocket.
type Topic struct {
	subscribers        map[*websocket.Conn]Subscriber
	onSubscribeMessage any
	active             bool
	activeMu           *sync.Mutex
	c                  chan any
}

// NewTopic creates a new [Topic].
func NewTopic(onSubscribMessage any) *Topic {
	return &Topic{
		subscribers:        make(map[*websocket.Conn]Subscriber),
		onSubscribeMessage: onSubscribMessage,
		active:             false,
		activeMu:           &sync.Mutex{},
		c:                  make(chan any, 100), //nolint:mnd //no magic number
	}
}

// Subscribe subscribes a [Subscriber] to this [Topic].
// If configured a message will be sent on subscribing.
// If no message handling go routine was
// running this will be started now.
func (t *Topic) Subscribe(ctx context.Context, conn *websocket.Conn) {
	sub := Subscriber{
		ctx:   context.WithoutCancel(ctx),
		topic: t,
		conn:  conn,
	}

	t.subscribers[conn] = sub

	if t.onSubscribeMessage != nil {
		sub.SendMessage(t.onSubscribeMessage)
	}

	if !t.active && t.activeMu.TryLock() {
		t.active = true

		go t.handleMessages()

		t.activeMu.Unlock()
	}
}

// UnSubscribe unsubscribes a [Subscriber] from this [Topic].
func (t *Topic) UnSubscribe(sub Subscriber) {
	delete(t.subscribers, sub.conn)
}

// EnqueueMessage enqueues a message if there are subscribers on this [Topic].
func (t *Topic) EnqueueMessage(msg any) {
	if !t.active {
		return
	}

	t.c <- msg
}

func (t *Topic) handleMessages() {
	for len(t.subscribers) > 0 {
		msg := <-t.c

		for _, sub := range t.subscribers {
			sub.SendMessage(msg)
		}
	}

	t.active = false
}
