package ws_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	wstools "github.com/xdoubleu/essentia/pkg/communication/ws"
	errortools "github.com/xdoubleu/essentia/pkg/errors"
	"github.com/xdoubleu/essentia/pkg/test"
	"github.com/xdoubleu/essentia/pkg/validate"
)

type TestResponse struct {
	Ok bool `json:"ok"`
}

type TestSubscribeMsg struct {
	TopicName string `json:"topicName"`
}

func (s TestSubscribeMsg) Validate() *validate.Validator {
	return validate.New()
}

func (s TestSubscribeMsg) Topic() string {
	return s.TopicName
}

func setup(t *testing.T) http.Handler {
	t.Helper()

	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](
		1,
		10,
	)

	_, err := ws.AddTopic(
		"exists",
		[]string{"http://localhost"},
		func(_ context.Context, _ *wstools.Topic) (any, error) { return TestResponse{Ok: true}, nil },
	)
	require.Nil(t, err)

	return ws.Handler()
}

func TestWebSocketExistingTopic(t *testing.T) {
	wsHandler := setup(t)

	tWeb := test.CreateWebSocketTester(wsHandler)
	tWeb.SetInitialMessage(TestSubscribeMsg{TopicName: "exists"})

	var initialResponse TestResponse
	err := tWeb.Do(t, &initialResponse, nil)

	require.Nil(t, err)
	assert.True(t, initialResponse.Ok)
}

func TestWebSocketUnknownTopic(t *testing.T) {
	wsHandler := setup(t)

	tWeb := test.CreateWebSocketTester(wsHandler)
	tWeb.SetInitialMessage(TestSubscribeMsg{TopicName: "unknown"})

	var initialResponse errortools.ErrorDto
	err := tWeb.Do(t, &initialResponse, nil)

	require.Nil(t, err)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), initialResponse.Error)
	assert.Equal(t, "topic 'unknown' doesn't exist", initialResponse.Message)
}

func TestWebSocketExistingHandler(t *testing.T) {
	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](1, 10)
	topic, err := ws.AddTopic(
		"exists",
		[]string{"http://localhost"},
		nil,
	)
	require.NotNil(t, topic)
	require.Nil(t, err)

	topic, err = ws.AddTopic(
		"exists",
		[]string{"http://localhost"},
		nil,
	)
	assert.Nil(t, topic)
	assert.EqualError(t, err, "topic 'exists' has already been added")
}
