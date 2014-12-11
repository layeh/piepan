package piepan

import (
	"time"

	"github.com/robertkrimen/otto"
)

type timer struct {
	cancel chan bool
}

func (t *timer) Cancel() {
	if t.cancel != nil {
		t.cancel <- true
	}
}

func (in *Instance) apiTimerNew(call otto.FunctionCall) otto.Value {
	callback := call.Argument(0)
	timeout, _ := call.Argument(1).ToInteger()

	t := &timer{
		cancel: make(chan bool),
	}

	go func() {
		defer func() {
			close(t.cancel)
			t.cancel = nil
		}()

		select {
		case <-time.After(time.Millisecond * time.Duration(timeout)):
			in.callValue(callback)
		case <-t.cancel:
		}
	}()

	ret, _ := in.state.ToValue(t)
	return ret
}
