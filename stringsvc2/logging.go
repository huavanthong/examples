package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

// TransportLoggingS3: Initialize a object from structure include log, and handler
type loggingMiddleware struct {
	logger log.Logger
	next   StringService
}

// TransportLoggingS4: Create handler for Uppercase
func (mw loggingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "uppercase",
			"input", s,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	// TransportLoggingS5: Use handler inside struct at here
	output, err = mw.next.Uppercase(s)
	return
}

func (mw loggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "count",
			"input", s,
			"n", n,
			"took", time.Since(begin),
		)
	}(time.Now())

	n = mw.next.Count(s)
	return
}
