package ratelimit

import (
	"time"

	"golang.org/x/time/rate"
)

var Limiter *rate.Limiter

func Init(workers int) {
	burst := workers / 2
	// A burst of 0 makes the underlying limiter reject every Wait()
	// call immediately instead of throttling, which silently disables
	// rate limiting for low worker counts (e.g. -w 1). Keep a floor of 1.
	if burst < 1 {
		burst = 1
	}
	Limiter = rate.NewLimiter(rate.Every(time.Millisecond), burst)
}
