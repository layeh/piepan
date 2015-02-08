package plugin

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/aarzilli/golua/lua"
	"github.com/layeh/piepan"
	"github.com/stevedonovan/luar"
)

func init() {
	piepan.Register("lua", &piepan.Plugin{
		Name: "Lua",
		New: func(in *piepan.Instance) piepan.Environment {
			s := luar.Init()
			p := &Plugin{
				instance:  in,
				state:     s,
				listeners: make(map[string][]*luar.LuaObject),
			}

			luar.Register(s, "piepan", luar.Map{
				"On":         p.apiOn,
				"Disconnect": p.apiDisconnect,
			})
			s.GetGlobal("piepan")
			s.NewTable()
			luar.Register(s, "*", luar.Map{
				"Play":       p.apiAudioPlay,
				"IsPlaying":  p.apiAudioIsPlaying,
				"Stop":       p.apiAudioStop,
				"NewTarget":  p.apiAudioNewTarget,
				"SetTarget":  p.apiAudioSetTarget,
				"Bitrate":    p.apiAudioBitrate,
				"SetBitrate": p.apiAudioSetBitrate,
				"Volume":     p.apiAudioVolume,
				"SetVolume":  p.apiAudioSetVolume,
			})
			s.SetField(-2, "Audio")
			s.NewTable()
			luar.Register(s, "*", luar.Map{
				"New": p.apiTimerNew,
			})
			s.SetField(-2, "Timer")
			s.NewTable()
			luar.Register(s, "*", luar.Map{
				"New": p.apiProcessNew,
			})
			s.SetField(-2, "Process")
			s.SetTop(0)

			return p
		},
	})
}

type Plugin struct {
	instance *piepan.Instance

	stateLock sync.Mutex
	state     *lua.State

	listeners map[string][]*luar.LuaObject
}

func (p *Plugin) LoadScriptFile(filename string) error {
	return p.state.DoFile(filename)
}

func (p *Plugin) apiOn(l *lua.State) int {
	event := strings.ToLower(l.CheckString(1))
	function := luar.NewLuaObject(l, 2)
	p.listeners[event] = append(p.listeners[event], function)
	return 0
}

func (p *Plugin) apiDisconnect(l *lua.State) int {
	if client := p.instance.Client; client != nil {
		client.Disconnect()
	}
	return 0
}

func (p *Plugin) error(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
}

func (p *Plugin) callValue(callback *luar.LuaObject, args ...interface{}) {
	p.stateLock.Lock()
	if _, err := callback.Call(args...); err != nil {
		p.error(err)
	}
	p.stateLock.Unlock()
}
