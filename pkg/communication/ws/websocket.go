package ws

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xdoubleu/essentia/pkg/validate"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// SubscribeMessageDto is implemented by all messages
// used to subscribe to a certain handler of a [WebSocketHandler].
type SubscribeMessageDto interface {
	validate.ValidatedType
	Topic() string
}

// A WebSocketHandler handles incoming requests to a
// websocket and makes sure subscriptions are made to the right topics.
type WebSocketHandler[T SubscribeMessageDto] struct {
	maxTopicWorkers        int
	topicChannelBufferSize int
	allowedOrigins         []string
	topicMap               map[string]*Topic
}

// CreateWebSocketHandler creates a new [WebSocketHandler].
func CreateWebSocketHandler[T SubscribeMessageDto](
	maxTopicWorkers int,
	topicChannelBufferSize int,
	allowedOrigins []string,
) WebSocketHandler[T] {
	for i, url := range allowedOrigins {
		if strings.Contains(url, "://") {
			allowedOrigins[i] = strings.Split(url, "://")[1]
		}
	}

	return WebSocketHandler[T]{
		maxTopicWorkers:        maxTopicWorkers,
		topicChannelBufferSize: topicChannelBufferSize,
		allowedOrigins:         allowedOrigins,
		topicMap:               make(map[string]*Topic),
	}
}

// AddTopic adds a topic to which can be subscribed using a [SubscribeMessageDto].
// The onSubscribeCallback is called for each
// new subscriber to fetch data to send them back.
func (h *WebSocketHandler[T]) AddTopic(
	topicName string,
	onSubscribeCallback OnSubscribeCallback,
) (*Topic, error) {
	_, ok := h.topicMap[topicName]
	if ok {
		return nil, fmt.Errorf("topic '%s' has already been added", topicName)
	}

	topic := NewTopic(
		topicName,
		h.maxTopicWorkers,
		h.topicChannelBufferSize,
		onSubscribeCallback,
	)
	h.topicMap[topicName] = topic

	return topic, nil
}

// UpdateTopicName updates the name of a topic without losing its subscribers.
func (h *WebSocketHandler[T]) UpdateTopicName(
	topic *Topic,
	newName string,
) (*Topic, error) {
	newTopic, ok := h.topicMap[topic.Name]
	if !ok {
		return nil, fmt.Errorf("topic '%s' doesn't exist", topic.Name)
	}

	_, ok = h.topicMap[newName]
	if ok {
		return nil, fmt.Errorf("topic '%s' already exists", newName)
	}

	newTopic.Name = newName
	delete(h.topicMap, topic.Name)
	h.topicMap[newTopic.Name] = newTopic

	return newTopic, nil
}

// RemoveTopic removes a topic to which can be subscribed using a [SubscribeMessageDto].
func (h *WebSocketHandler[T]) RemoveTopic(topic *Topic) error {
	_, ok := h.topicMap[topic.Name]
	if !ok {
		return fmt.Errorf("topic '%s' doesn't exist", topic.Name)
	}

	topic.StopPool()
	delete(h.topicMap, topic.Name)
	return nil
}

// Handler returns the [http.HandlerFunc] of a [WebSocketHandler].
func (h WebSocketHandler[T]) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//nolint:exhaustruct //other fields are optional
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			//todo OriginPatterns: h.allowedOrigins,
			InsecureSkipVerify: true,
		})
		if err != nil {
			UpgradeErrorResponse(w, r, err)
			return
		}

		// in case you want to subscribe on multiple topics
		for {
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

			topic, ok := h.topicMap[msg.Topic()]
			if !ok {
				ErrorResponse(
					r.Context(),
					conn,
					http.StatusBadRequest,
					fmt.Sprintf("topic '%s' doesn't exist", msg.Topic()),
				)
				return
			}

			err = topic.Subscribe(conn)
			if err != nil {
				ServerErrorResponse(r.Context(), conn, err)
				return
			}
		}
	}
}
