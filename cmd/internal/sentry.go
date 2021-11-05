package internal

import (
	"github.com/Kichiyaki/appmode"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"os"
)

func InitSentry(release string) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Environment:      appmode.Get(),
		Release:          release,
		Debug:            false,
		TracesSampleRate: 0.3,
	})
	if err != nil {
		return errors.Wrap(err, "sentry.Init")
	}

	return nil
}
