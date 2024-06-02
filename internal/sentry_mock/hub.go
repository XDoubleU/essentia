package sentry_mock

import "github.com/getsentry/sentry-go"

const testDsn = "http://whatever@example.com/1337"

func GetMockedClientOptions() *sentry.ClientOptions {
	return &sentry.ClientOptions{
		Dsn:       testDsn,
		Transport: &TransportMock{},
	}
}

func GetMockedHub() *sentry.Hub {
	clientOptions := GetMockedClientOptions()

	client, err := sentry.NewClient(*clientOptions)
	if err != nil {
		panic(err)
	}

	scope := sentry.NewScope()
	hub := sentry.NewHub(client, scope)
	return hub
}
