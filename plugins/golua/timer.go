package plugin

import (
	"time"

	"github.com/yuin/gopher-lua"
)

type Timer struct {
	cancel chan bool
}

func (t *Timer) Cancel() {
	if t.cancel != nil {
		t.cancel <- true
	}
}

func (p *Plugin) apiTimerNew(callback *lua.LFunction, timeout int) *Timer {
	t := &Timer{
		cancel: make(chan bool),
	}

	go func() {
		defer func() {
			close(t.cancel)
			t.cancel = nil
		}()

		select {
		case <-time.After(time.Millisecond * time.Duration(timeout)):
			p.callValue(callback)
		case <-t.cancel:
		}
	}()

	return t
}
