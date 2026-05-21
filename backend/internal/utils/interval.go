package utils

import "time"

func Interval(f func(), interval time.Duration) {
	go func() {
		f()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			f()
		}
	}()
}
