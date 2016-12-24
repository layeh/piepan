package piepan // import "layeh.com/piepan"

import (
	"github.com/yuin/gopher-lua"
)

func (s *State) apiDisconnect() {
	if client := s.Client; client != nil {
		client.Disconnect()
	}
}

func (s *State) apiOn(event string, fn *lua.LFunction) {
	s.listeners[event] = append(s.listeners[event], fn)
}
