package diplomat

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

func stringSliceToInterfaceSlice(slice []string) []interface{} {
	pointers := make([]interface{}, len(slice))
	for i, v := range slice {
		s := v
		pointers[i] = &s
	}
	return pointers
}

func throttle(interval time.Duration, events <-chan fsnotify.Event) <-chan fsnotify.Event {
	c := make(chan fsnotify.Event)
	go func() {
		var lastEvent *fsnotify.Event
		var ticker *time.Ticker
		for {
			if ticker == nil {
				ticker = time.NewTicker(interval)
			}
			select {
			case <-ticker.C:
				if lastEvent != nil {
					c <- *lastEvent
					lastEvent = nil
				}
				break
			case e := <-events:
				lastEvent = &e
				ticker = nil
			}
		}
	}()
	return c
}
