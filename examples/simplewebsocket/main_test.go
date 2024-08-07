package main

import (
	"io"
	"log"
	"testing"

	"github.com/XDoubleU/essentia/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebsocket(t *testing.T) {
	app := NewApp(log.New(io.Discard, "", 0))
	app.config.Env = TestEnv

	routes, err := app.Routes()
	require.Nil(t, err)

	tWeb := test.CreateWebsocketTester(*routes)
	tWeb.SetInitialMessage(SubjectMessageDto{
		Subject: "subject",
	})

	var initialResponse ResponseMessageDto
	err = tWeb.Do(t, &initialResponse, nil)
	require.Nil(t, err)

	assert.Equal(t, "Hello, World!", initialResponse.Message)
}
