package piepan

import (
	"fmt"
	"os"
	"time"

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type timer struct {
	cancel chan bool
}

func (t *timer) Cancel() {
	if t.cancel != nil {
		t.cancel <- true
	}
}

func (in *Instance) timerNew(l *lua.State) int {
	callback := luar.NewLuaObject(l, 1)
	timeout := l.CheckInteger(2)

	t := &timer{
		cancel: make(chan bool),
	}

	go func() {
		defer callback.Close()
		defer func() {
			close(t.cancel)
			t.cancel = nil
		}()

		select {
		case <-time.After(time.Second * time.Duration(timeout)):
			in.stateLock.Lock()
			defer in.stateLock.Unlock()

			if _, err := callback.Call(); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
		case <-t.cancel:
		}
	}()

	obj := luar.NewLuaObjectFromValue(l, t)
	obj.Push()
	return 1
}
