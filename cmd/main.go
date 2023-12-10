package main

import (
	"fmt"
	"time"

	"github.com/evebot-tools/#REPO#/internal"
	"github.com/evebot-tools/libs/utils"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              utils.GetEnv("SENTRY_DSN", ""),
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	defer sentry.Flush(5 * time.Second)
	internal.Run()
}
