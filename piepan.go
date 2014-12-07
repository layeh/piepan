package piepan

import (
	"strings"
	"sync"

	"github.com/aarzilli/golua/lua"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/stevedonovan/luar"
)

type Instance struct {
	client *gumble.Client

	audio             *gumble_ffmpeg.Stream
	audioCallbackFunc *luar.LuaObject

	stateLock sync.Mutex
	state     *lua.State
	listeners map[string][]*luar.LuaObject
}

func New(client *gumble.Client) *Instance {
	instance := &Instance{
		client:    client,
		state:     luar.Init(),
		listeners: make(map[string][]*luar.LuaObject),
	}
	instance.audio, _ = gumble_ffmpeg.New(instance.client)
	instance.audio.Done = instance.audioCallback

	luar.Register(instance.state, "piepan", luar.Map{
		"On": instance.luaPiepanOn,
	})

	instance.state.GetGlobal("piepan")

	instance.state.NewTable()
	luar.Register(instance.state, "*", luar.Map{
		"Play":        instance.audioPlay,
		"IsPlaying":   instance.audioIsPlaying,
		"Stop":        instance.audioStop,
		"SetCallback": instance.audioSetCallback,
	})
	instance.state.SetField(-2, "Audio")

	instance.state.NewTable()
	luar.Register(instance.state, "*", luar.Map{
		"New": instance.timerNew,
	})
	instance.state.SetField(-2, "Timer")

	instance.state.NewTable()
	luar.Register(instance.state, "*", luar.Map{
		"New": instance.processNew,
	})
	instance.state.SetField(-2, "Process")

	instance.state.SetTop(0)
	return instance
}

func (in *Instance) luaPiepanOn(l *lua.State) int {
	event := strings.ToLower(l.CheckString(1))
	function := luar.NewLuaObject(l, 2)
	in.listeners[event] = append(in.listeners[event], function)
	return 0
}

func (in *Instance) LoadScriptFile(filename string) error {
	return in.state.DoFile(filename)
}

func (in *Instance) Destroy() {
	in.state.Close()
}
