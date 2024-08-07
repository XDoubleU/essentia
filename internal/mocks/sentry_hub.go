package mocks

import "github.com/getsentry/sentry-go"

const testDsn = "http://whatever@example.com/1337"

func GetMockedSentryClientOptions() *sentry.ClientOptions {
	//nolint:exhaustruct //other fields are optional
	return &sentry.ClientOptions{
		Dsn:       testDsn,
		Transport: NewTransportMock(),
	}
}

func GetMockedSentryHub() *sentry.Hub {
	clientOptions := GetMockedSentryClientOptions()

	client, err := sentry.NewClient(*clientOptions)
	if err != nil {
		panic(err)
	}

	scope := sentry.NewScope()
	hub := sentry.NewHub(client, scope)
	return hub
}
