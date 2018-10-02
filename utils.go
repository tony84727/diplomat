package diplomat

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

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

func newFanoutHub() *fanoutHub {
	source := make(chan interface{}, 10)
	return &fanoutHub{
		source:     source,
		listeners:  make([]chan<- interface{}, 0, 1),
		addingChan: make(chan chan<- interface{}),
	}
}

// fanoutHub send received message to its listeners(channels)
type fanoutHub struct {
	source     chan interface{}
	listeners  []chan<- interface{}
	addingChan chan chan<- interface{}
}

func (fh *fanoutHub) addListener(c chan<- interface{}) {
	fh.addingChan <- c
}

func (fh *fanoutHub) broadcast(message interface{}) {
	fh.source <- message
}

func (fh *fanoutHub) run() {
	for {
		select {
		case newListener := <-fh.addingChan:
			fh.listeners = append(fh.listeners, newListener)
			break
		case m, ok := <-fh.source:
			if !ok {
				fh.source = nil
				break
			}
			var wg sync.WaitGroup
			for _, l := range fh.listeners {
				wg.Add(1)
				go fh.forwardMessage(l, m, &wg)
			}
			wg.Wait()
		}
		if fh.source == nil {
			break
		}
	}
	for _, l := range fh.listeners {
		close(l)
	}
}

func (fh fanoutHub) forwardMessage(channel chan<- interface{}, message interface{}, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Millisecond)
	select {
	case channel <- message:
		break
	case <-ticker.C:
		break
	}
	wg.Done()
}

type pathBasedThrottler struct {
	threshold time.Duration
	channels  map[string]chan<- fsnotify.Event
}

func newPathBasedThrottler(threshold time.Duration) *pathBasedThrottler {
	return &pathBasedThrottler{
		threshold: threshold,
		channels:  make(map[string]chan<- fsnotify.Event),
	}
}

func (pbt *pathBasedThrottler) run(input <-chan fsnotify.Event) <-chan fsnotify.Event {
	output := make(chan fsnotify.Event)
	// distribute
	go func() {
		for i := range input {
			if _, ok := pbt.channels[i.Name]; !ok {
				throttled := make(chan fsnotify.Event)
				go func() {
					for e := range throttle(pbt.threshold, throttled) {
						output <- e
					}
				}()
				pbt.channels[i.Name] = throttled
			}
			pbt.channels[i.Name] <- i
		}
	}()
	return output
}
