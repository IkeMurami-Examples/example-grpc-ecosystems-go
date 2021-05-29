package utils

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

const sentryFlushTime = time.Second * 5

//SentrySetup - setup sentry error logger
func SentrySetup(dsn string) {
	sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
	})
	sentry.Flush(sentryFlushTime)
	zerolog.ErrorMarshalFunc = func(err error) interface{} {
		sentry.CaptureException(err)
		return err
	}
}
