package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xdoubleu/essentia/pkg/config"
	"github.com/xdoubleu/essentia/pkg/logging"
	"github.com/xdoubleu/essentia/pkg/test"
)

func TestWebSocket(t *testing.T) {
	app := NewApp(logging.NewNopLogger())
	app.config.Env = config.TestEnv

	tWeb := test.CreateWebSocketTester(app.Routes())
	tWeb.SetInitialMessage(SubscribeMessageDto{
		TopicName: "topic",
	})

	var initialResponse ResponseMessageDto
	err := tWeb.Do(t, &initialResponse, nil)
	require.Nil(t, err)

	assert.Equal(t, "Hello, World!", initialResponse.Message)
}
