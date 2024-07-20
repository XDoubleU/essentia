package wstools

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/XDoubleU/essentia/pkg/validate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// SubscribeMessageDto is implemented by all messages
// used to subscribe to a certain handler of a [WebSocketHandler].
type SubscribeMessageDto interface {
	validate.ValidatedType
	GetTopicName() string
}

// A WebSocketHandler handles incoming requests to a
// websocket and makes sure subscriptions are made to the right topics.
type WebSocketHandler[T SubscribeMessageDto] struct {
	allowedOrigins []string
	topicMap       map[string]*Topic
}

// CreateWebSocketHandler creates a new [WebSocketHandler].
func CreateWebSocketHandler[T SubscribeMessageDto](
	allowedOrigins []string,
) WebSocketHandler[T] {
	for i, url := range allowedOrigins {
		if strings.Contains(url, "://") {
			allowedOrigins[i] = strings.Split(url, "://")[1]
		}
	}

	return WebSocketHandler[T]{
		allowedOrigins: allowedOrigins,
		topicMap:       make(map[string]*Topic),
	}
}

// AddTopic adds a topic to which can be subscribed using a [SubscribeMessageDto].
// The optional onSubscribeMessage is sent to each new subscriber.
func (h *WebSocketHandler[T]) AddTopic(
	topicName string,
	onSubscribeMessage any,
) (*Topic, error) {
	_, ok := h.topicMap[topicName]
	if ok {
		return nil, fmt.Errorf("topic '%s' has already been added", topicName)
	}

	topic := NewTopic(onSubscribeMessage)
	h.topicMap[topicName] = topic

	return topic, nil
}

// Handler returns the [http.HandlerFunc] of a [WebSocketHandler].
func (h WebSocketHandler[T]) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//nolint:exhaustruct //other fields are optional
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			OriginPatterns: h.allowedOrigins,
		})
		if err != nil {
			UpgradeErrorResponse(w, r, err)
			return
		}

		var msg T
		err = wsjson.Read(r.Context(), conn, &msg)
		if err != nil {
			ServerErrorResponse(r.Context(), conn, err)
			return
		}

		if v := msg.Validate(); !v.Valid() {
			FailedValidationResponse(r.Context(), conn, v.Errors)
			return
		}

		topic, ok := h.topicMap[msg.GetTopicName()]
		if !ok {
			ErrorResponse(
				r.Context(),
				conn,
				http.StatusBadRequest,
				fmt.Sprintf("topic '%s' doesn't exist", msg.GetTopicName()),
			)
			return
		}

		topic.Subscribe(r.Context(), conn)
	}
}
