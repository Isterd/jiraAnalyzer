// pkg/backoff/backoff.go
package backoff

import (
	"math"
	"time"
)

type BackOff struct {
	MaxRetries   int
	MinTimeSleep time.Duration
	MaxTimeSleep time.Duration
}

func (b *BackOff) Duration(attempt int) time.Duration {
	if attempt >= b.MaxRetries {
		return 0
	}
	return time.Duration(math.Min(
		float64(b.MinTimeSleep)*math.Pow(2, float64(attempt)),
		float64(b.MaxTimeSleep),
	)) * time.Millisecond
}

func NewExponentialBackOff(min, max time.Duration) BackOff {
	return BackOff{
		MinTimeSleep: min * time.Millisecond,
		MaxTimeSleep: max * time.Millisecond,
		MaxRetries:   5,
	}
}
