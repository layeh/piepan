package plugin

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/layeh/gopher-luar"
	"github.com/layeh/piepan"
	"github.com/yuin/gopher-lua"
)

func init() {
	piepan.Register("lua", &piepan.Plugin{
		Name: "Lua (Go)",
		New: func(in *piepan.Instance) piepan.Environment {
			s := lua.NewState()
			p := &Plugin{
				instance:  in,
				state:     s,
				listeners: make(map[string][]lua.LValue),
			}

			pp := s.NewTable()
			s.SetGlobal("piepan", pp)

			pp.RawSetH(lua.LString("On"), luar.New(s, p.apiOn))
			pp.RawSetH(lua.LString("Disconnect"), luar.New(s, p.apiDisconnect))

			t := s.NewTable()
			t.RawSetH(lua.LString("Play"), luar.New(s, p.apiAudioPlay))
			t.RawSetH(lua.LString("IsPlaying"), luar.New(s, p.apiAudioIsPlaying))
			t.RawSetH(lua.LString("Stop"), luar.New(s, p.apiAudioStop))
			t.RawSetH(lua.LString("NewTarget"), luar.New(s, p.apiAudioNewTarget))
			t.RawSetH(lua.LString("SetTarget"), luar.New(s, p.apiAudioSetTarget))
			t.RawSetH(lua.LString("Bitrate"), luar.New(s, p.apiAudioBitrate))
			t.RawSetH(lua.LString("SetBitrate"), luar.New(s, p.apiAudioSetBitrate))
			t.RawSetH(lua.LString("Volume"), luar.New(s, p.apiAudioVolume))
			t.RawSetH(lua.LString("SetVolume"), luar.New(s, p.apiAudioSetVolume))
			pp.RawSetH(lua.LString("Audio"), t)

			t = s.NewTable()
			t.RawSetH(lua.LString("New"), luar.New(s, p.apiTimerNew))
			pp.RawSetH(lua.LString("Timer"), t)

			t = s.NewTable()
			t.RawSetH(lua.LString("New"), luar.New(s, p.apiProcessNew))
			pp.RawSetH(lua.LString("Process"), t)

			return p
		},
	})
}

type Plugin struct {
	instance *piepan.Instance

	stateLock sync.Mutex
	state     *lua.LState

	listeners map[string][]lua.LValue
}

func (p *Plugin) LoadScriptFile(filename string) error {
	return p.state.DoFile(filename)
}

func (p *Plugin) apiOn(event string, fn *lua.LFunction) {
	p.listeners[event] = append(p.listeners[event], fn)
}

func (p *Plugin) apiDisconnect() {
	if client := p.instance.Client; client != nil {
		client.Disconnect()
	}
}

func (p *Plugin) error(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
}

func (p *Plugin) callValue(callback lua.LValue, args ...interface{}) {
	p.stateLock.Lock()
	p.state.Push(callback)
	for _, arg := range args {
		p.state.Push(luar.New(p.state, arg))
	}
	p.state.PCall(len(args), 0, p.state.NewFunction(func(L *lua.LState) int {
		p.error(errors.New(L.CheckString(1)))
		return 0
	}))
	p.state.SetTop(0)
	p.stateLock.Unlock()
}
