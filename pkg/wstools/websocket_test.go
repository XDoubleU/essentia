package wstools_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/httptools"
	"github.com/xdoubleu/essentia/pkg/test"
	"github.com/xdoubleu/essentia/pkg/validate"
	"github.com/xdoubleu/essentia/pkg/wstools"
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

func (s TestSubscribeMsg) GetTopicName() string {
	return s.TopicName
}

func setup(t *testing.T) http.Handler {
	t.Helper()

	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](
		1,
		10,
		[]string{"http://localhost"},
	)

	_, err := ws.AddTopic(
		"exists",
		TestResponse{Ok: true},
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

	var initialResponse httptools.ErrorDto
	err := tWeb.Do(t, &initialResponse, nil)

	require.Nil(t, err)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), initialResponse.Error)
	assert.Equal(t, "topic 'unknown' doesn't exist", initialResponse.Message)
}

func TestWebSocketExistingHandler(t *testing.T) {
	ws := wstools.CreateWebSocketHandler[TestSubscribeMsg](1, 10, []string{"localhost"})
	topic, err := ws.AddTopic(
		"exists",
		nil,
	)
	require.NotNil(t, topic)
	require.Nil(t, err)

	topic, err = ws.AddTopic(
		"exists",
		nil,
	)
	assert.Nil(t, topic)
	assert.EqualError(t, err, "topic 'exists' has already been added")
}