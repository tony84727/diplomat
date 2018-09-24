package diplomat

import (
	"testing"
	"time"
)

func assertReceiveEvent(t *testing.T, listener <-chan interface{}, timeout time.Duration) {
	ticker := time.NewTicker(timeout)
	select {
	case <-listener:
		break
	case <-ticker.C:
		t.Log("didn't receive change notification")
		t.Fail()
		break
	}
}

func TestMaintainOutlineBroadcastChange(t *testing.T) {
	d := New()
	fakeOutline := &Outline{}
	l := make(chan interface{})
	d.changeListeners.addListener(l)
	d.SetOutline(fakeOutline)
	assertReceiveEvent(t, l, time.Second)
}

func TestMaintainTranslationsBroadcastChange(t *testing.T) {
	d := New()
	fakeTranslations := &PartialTranslation{}
	l := make(chan interface{})
	d.changeListeners.addListener(l)
	d.SetTranslation(fakeTranslations)
	assertReceiveEvent(t, l, time.Second)
}
