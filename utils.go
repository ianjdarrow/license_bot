package main

import (
	"time"
)

func getTimeSinceMs(t time.Time) time.Duration {
	return time.Now().Sub(t).Truncate(time.Millisecond)
}
