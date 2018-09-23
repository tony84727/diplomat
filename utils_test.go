package diplomat

import (
	"math"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

func fakeFsnotifyEvent(name string) fsnotify.Event {
	return fsnotify.Event{
		Name: name,
	}
}

type stopWatch struct {
	stamp time.Time
}

func (s *stopWatch) start() {
	s.stamp = time.Now()
}

func (s stopWatch) stop() time.Duration {
	return time.Now().Sub(s.stamp)
}

func TestThrottle(t *testing.T) {
	in := make(chan fsnotify.Event, 10)
	var timer stopWatch
	timer.start()
	in <- fakeFsnotifyEvent("first")
	in <- fakeFsnotifyEvent("second")
	expectedDuration := timer.stop() + time.Millisecond
	timer.start()
	out := throttle(expectedDuration, in)
	e := <-out
	actualInterval := timer.stop()
	if actualInterval <= expectedDuration {
		t.Logf(
			"expect at least %f, got %f, diff: %f",
			float64(expectedDuration/time.Millisecond),
			float64(actualInterval/time.Millisecond),
			math.Abs(float64(actualInterval-expectedDuration)/float64(time.Millisecond)),
		)
		t.Fail()
	}
	assert.Equal(t, "second", e.Name)
}
