package ratelimit

import (
	"time"

	"golang.org/x/time/rate"
)

var Limiter = rate.NewLimiter(rate.Every(time.Millisecond), 500)
