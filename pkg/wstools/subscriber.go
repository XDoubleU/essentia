package wstools

import (
	"context"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Subscriber is used to receive messages
// from a [Topic] and managed the [websocket.Conn].
type Subscriber struct {
	ctx   context.Context
	topic *Topic
	conn  *websocket.Conn
}

// SendMessage sends a message to this [Subscriber].
// If the connection would be closed,
// [UnSubscribe] will be called.
func (sub Subscriber) SendMessage(msg any) {
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
