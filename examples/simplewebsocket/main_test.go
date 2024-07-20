package main

import (
	"io"
	"log"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSocket(t *testing.T) {
	app := NewApp(log.New(io.Discard, "", 0))
	app.config.Env = TestEnv

	routes, err := app.Routes()
	require.Nil(t, err)

	tWeb := test.CreateWebSocketTester(*routes)
	tWeb.SetInitialMessage(SubscribeMessageDto{
		TopicName: "topic",
	})

	var initialResponse ResponseMessageDto
	err = tWeb.Do(t, &initialResponse, nil)
	require.Nil(t, err)

	assert.Equal(t, "Hello, World!", initialResponse.Message)
}
