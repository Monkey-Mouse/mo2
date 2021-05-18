package mo2ticker

import (
	"time"
)

func ExecuteFunc(duration time.Duration, handler func()) chan struct{} {
	ticker := time.NewTicker(duration)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				ticker.Stop()
				return
			case _ = <-ticker.C:
				handler()
			}
		}
	}()
	return stop
}
