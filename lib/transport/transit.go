package transport

import (
	"time"
)

func MutableTransit[SrcV any, DstV any](src chan SrcV, dst chan DstV, mutator func(srcValue SrcV) DstV, delay time.Duration) {
	var ticker *time.Ticker

	if delay > 0 {
		ticker = time.NewTicker(delay)
	}

	for v := range src {
		dst <- mutator(v)

		if ticker != nil {
			<-ticker.C
		}
	}

	if ticker != nil {
		ticker.Stop()
	}
}
