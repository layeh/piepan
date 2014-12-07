package piepan

import (
	"fmt"
	"os"

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

func (in *Instance) audioPlay(l *lua.State) int {
	if in.audio.IsPlaying() {
		return 0
	}

	obj := luar.NewLuaObject(l, 1)
	defer obj.Close()

	filename := obj.Get("filename").(string)
	callback := obj.GetObject("callback")

	if callback.Type != "nil" {
		in.audio.Done = func() {
			in.stateLock.Lock()
			defer in.stateLock.Unlock()

			if _, err := callback.Call(); err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}

			callback.Close()
		}
	}

	in.audio.Play(filename)
	return 0
}

func (in *Instance) audioStop() {
	in.audio.Stop()
}

func (in *Instance) audioIsPlaying() bool {
	return in.audio.IsPlaying()
}
