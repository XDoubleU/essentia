package wstools

import (
	"context"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Subscriber is used to receive messages
// from a [Topic] and managed the [websocket.Conn].
type Subscriber struct {
	id    string
	ctx   context.Context
	topic *Topic
	conn  *websocket.Conn
}

func NewSubscriber(ctx context.Context, topic *Topic, conn *websocket.Conn) Subscriber {
	return Subscriber{
		id:    uuid.NewString(),
		ctx:   ctx,
		topic: topic,
		conn:  conn,
	}
}

func (sub Subscriber) ID() string {
	return sub.id
}

// ExecuteCallback sends a message to this [Subscriber].
// If the connection would be closed,
// [UnSubscribe] will be called.
func (sub Subscriber) ExecuteCallback(msg any) {
	err := wsjson.Write(sub.ctx, sub.conn, msg)
	if err == nil {
		return
	}

	if websocket.CloseStatus(err) != -1 {
		sub.topic.UnSubscribe(sub)
		return
	}

	ServerErrorResponse(sub.ctx, sub.conn, err)
}
