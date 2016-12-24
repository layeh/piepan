package piepan // import "layeh.com/piepan"

import (
	"time"

	"github.com/yuin/gopher-lua"
)

type timer struct {
	cancel chan struct{}
}

func (t *timer) Cancel() {
	if t.cancel != nil {
		t.cancel <- struct{}{}
	}
}

func (s *State) apiTimerNew(callback *lua.LFunction, timeout int) *timer {
	t := &timer{
		cancel: make(chan struct{}),
	}

	go func() {
		defer func() {
			close(t.cancel)
			t.cancel = nil
		}()

		select {
		case <-time.After(time.Millisecond * time.Duration(timeout)):
			s.callValue(callback)
		case <-t.cancel:
		}
	}()

	return t
}
