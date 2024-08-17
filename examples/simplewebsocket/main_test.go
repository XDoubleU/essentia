package main

import (
	"testing"

	"github.com/XDoubleU/essentia/pkg/config"
	"github.com/XDoubleU/essentia/pkg/logging"
	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSocket(t *testing.T) {
	cfg := NewConfig()
	cfg.Env = config.TestEnv

	app := NewApp(logging.NewNopLogger(), cfg)

	tWeb := test.CreateWebSocketTester(app.Routes())
	tWeb.SetInitialMessage(SubscribeMessageDto{
		TopicName: "topic",
	})

	var initialResponse ResponseMessageDto
	err := tWeb.Do(t, &initialResponse, nil)
	require.Nil(t, err)

	assert.Equal(t, "Hello, World!", initialResponse.Message)
}
