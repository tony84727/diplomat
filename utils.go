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
		tickChan := make(chan time.Time)
		for {
			select {
			case e, ok := <-events:
				if !ok {
					events = nil
				}
				lastEvent = &e
				if ticker != nil {
					ticker.Stop()
				}
				ticker = time.NewTicker(interval)
				go func() {
					for t := range ticker.C {
						tickChan <- t
					}
				}()
				break
			case <-tickChan:
				ticker.Stop()
				ticker = nil
				if lastEvent != nil {
					c <- *lastEvent
				}
			}
			if events == nil {
				break
			}
		}
	}()
	return c
}
