package main

import (
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/test"
)

func TestWebSocket(t *testing.T) {
	app := NewApp(log.New(io.Discard, "", 0))
	app.config.Env = config.TestEnv

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
