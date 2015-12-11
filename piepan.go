package piepan

import (
	"fmt"
	"os"
	"sync"

	"github.com/layeh/gopher-luar"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleffmpeg"
	"github.com/yuin/gopher-lua"
)

type State struct {
	Client *gumble.Client

	LState *lua.LState
	table  *lua.LTable

	audioStream  *gumbleffmpeg.Stream
	audioVolume  float32
	AudioCommand string

	mu        sync.Mutex
	listeners map[string][]lua.LValue
}

func New(client *gumble.Client) *State {
	l := lua.NewState()
	state := &State{
		Client:    client,
		LState:    l,
		listeners: make(map[string][]lua.LValue),
	}
	t := l.NewTable()
	t.RawSetString("On", luar.New(l, state.apiOn))
	t.RawSetString("Disconnect", luar.New(l, state.apiDisconnect))
	state.table = t
	l.SetGlobal("piepan", t)
	{
		s := l.NewTable()
		s.RawSetString("Play", luar.New(l, state.apiAudioPlay))
		s.RawSetString("IsPlaying", luar.New(l, state.apiAudioIsPlaying))
		s.RawSetString("Stop", luar.New(l, state.apiAudioStop))
		s.RawSetString("NewTarget", luar.New(l, state.apiAudioNewTarget))
		s.RawSetString("SetTarget", luar.New(l, state.apiAudioSetTarget))
		s.RawSetString("Bitrate", luar.New(l, state.apiAudioBitrate))
		s.RawSetString("SetBitrate", luar.New(l, state.apiAudioSetBitrate))
		s.RawSetString("Volume", luar.New(l, state.apiAudioVolume))
		s.RawSetString("SetVolume", luar.New(l, state.apiAudioSetVolume))
		t.RawSetString("Audio", s)
	}
	{
		s := l.NewTable()
		s.RawSetString("New", luar.New(l, state.apiTimerNew))
		t.RawSetString("Timer", t)
	}
	{
		s := l.NewTable()
		s.RawSetString("New", luar.New(l, state.apiProcessNew))
		t.RawSetString("Process", s)
	}
	client.Attach(state)
	return state
}

func (s *State) LoadFile(filename string) error {
	return s.LState.DoFile(filename)
}

func (s *State) callValue(callback lua.LValue, args ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.LState.Push(callback)
	for _, arg := range args {
		s.LState.Push(luar.New(s.LState, arg))
	}
	s.LState.PCall(len(args), 0, s.LState.NewFunction(func(L *lua.LState) int {
		fmt.Fprintf(os.Stderr, "%s\n", L.CheckString(1))
		return 0
	}))
	s.LState.SetTop(0)
}
