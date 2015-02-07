package plugin

import (
	"time"

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type Timer struct {
	cancel chan bool
}

func (t *Timer) Cancel() {
	if t.cancel != nil {
		t.cancel <- true
	}
}

func (p *Plugin) apiTimerNew(l *lua.State) int {
	callback := luar.NewLuaObject(l, 1)
	timeout := l.ToInteger(2)

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
			callback.Close()
		case <-t.cancel:
		}
	}()

	obj := luar.NewLuaObjectFromValue(l, t)
	obj.Push()
	obj.Close()
	return 1
}
