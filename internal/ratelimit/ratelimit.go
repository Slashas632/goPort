package ratelimit

import (
	"time"

	"golang.org/x/time/rate"
)

var Limiter *rate.Limiter

func Init(workers int) {
	Limiter = rate.NewLimiter(rate.Every(time.Millisecond), workers/2)
}
