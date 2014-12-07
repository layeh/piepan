package piepan

import (
	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

func (in *Instance) audioPlay(filename string) {
	in.audio.Play(filename)
}

func (in *Instance) audioSetCallback(l *lua.State) int {
	if in.audioCallbackFunc != nil {
		in.audioCallbackFunc.Close()
	}
	in.audioCallbackFunc = luar.NewLuaObject(l, 1)
	return 0
}

func (in *Instance) audioCallback() {
	if callback := in.audioCallbackFunc; callback != nil {
		callback.Call()
	}
}

func (in *Instance) audioStop() {
	in.audio.Stop()
}

func (in *Instance) audioIsPlaying() bool {
	return in.audio.IsPlaying()
}
