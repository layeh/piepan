package piepan // import "layeh.com/piepan"

import (
	"sync"
	"time"

	"layeh.com/gopus"
	"layeh.com/gumble/gumble"
	"layeh.com/gumble/gumbleffmpeg"
	"layeh.com/gumble/opus"
	"github.com/yuin/gopher-lua"
)

type audioStream struct {
	state    *State
	s        *gumbleffmpeg.Stream
	callback lua.LValue
	wg       sync.WaitGroup
}

func (a *audioStream) Volume() float32 {
	return a.s.Volume
}

func (a *audioStream) SetVolume(volume float32) {
	a.s.Volume = volume
}

func (a *audioStream) IsPlaying() bool {
	return a.s.State() == gumbleffmpeg.StatePlaying
}

func (a *audioStream) Play() {
	a.state.streamMu.Lock()
	if a.state.stream != nil {
		a.state.streamMu.Unlock()
		return
	}
	a.state.stream = a
	a.state.streamMu.Unlock()
	a.wg.Add(1)
	err := a.s.Play()
	if err != nil {
		a.state.streamMu.Lock()
		a.state.stream = nil
		a.state.streamMu.Unlock()
		a.wg.Done()
		panic(err.Error())
	}
	go func() {
		for _, listener := range a.state.listeners["stream"] {
			a.state.callValue(listener, a)
		}

		a.s.Wait()
		a.state.streamMu.Lock()
		a.state.stream = nil
		a.state.streamMu.Unlock()
		a.wg.Done()

		for _, listener := range a.state.listeners["stream"] {
			a.state.callValue(listener, a)
		}

		if a.callback.Type() != lua.LTNil {
			a.state.callValue(a.callback)
		}
	}()
}

func (a *audioStream) IsStopped() bool {
	return a.s.State() == gumbleffmpeg.StateStopped
}

func (a *audioStream) Stop() {
	a.s.Stop()
	a.wg.Wait()
}

func (a *audioStream) IsPaused() bool {
	return a.s.State() == gumbleffmpeg.StatePaused
}

func (a *audioStream) Pause() {
	a.state.streamMu.Lock()
	defer a.state.streamMu.Unlock()
	if a.state.stream != a {
		return
	}
	a.state.stream = nil
	a.s.Pause()
	go func() {
		for _, listener := range a.state.listeners["stream"] {
			a.state.callValue(listener, a)
		}
	}()
}

func (a *audioStream) Elapsed() float64 {
	return a.s.Elapsed().Seconds()
}

func (s *State) apiAudioNew(tbl *lua.LTable) *audioStream {
	filename := tbl.RawGetString("filename")

	exec := tbl.RawGetString("exec")
	args := tbl.RawGetString("args")

	offset := tbl.RawGetString("offset")
	callback := tbl.RawGetString("callback")

	if enc, ok := s.Client.AudioEncoder.(*opus.Encoder); ok {
		enc.SetApplication(gopus.Audio)
	}

	var source gumbleffmpeg.Source
	switch {
	// source file
	case filename != lua.LNil && exec == lua.LNil && args == lua.LNil:
		source = gumbleffmpeg.SourceFile(filename.String())
	// source exec
	case filename == lua.LNil && exec != lua.LNil:
		var argsStr []string
		if argsTable, ok := args.(*lua.LTable); ok {
			for i := 1; ; i++ {
				arg := argsTable.RawGetInt(i)
				if arg == lua.LNil {
					break
				}
				argsStr = append(argsStr, arg.String())
			}
		}
		source = gumbleffmpeg.SourceExec(exec.String(), argsStr...)
	default:
		panic("invalid piepan.Audio.New source type")
	}

	stream := &audioStream{
		state:    s,
		s:        gumbleffmpeg.New(s.Client, source),
		callback: callback,
	}
	stream.s.Command = s.AudioCommand

	if number, ok := offset.(lua.LNumber); ok {
		stream.s.Offset = time.Second * time.Duration(number)
	}

	return stream
}

func (s *State) apiAudioNewTarget(id uint32) *gumble.VoiceTarget {
	return &gumble.VoiceTarget{
		ID: id,
	}
}

func (s *State) apiAudioBitrate() int {
	if enc, ok := s.Client.AudioEncoder.(*opus.Encoder); ok {
		return enc.Bitrate()
	}
	return -1
}

func (s *State) apiAudioSetBitrate(bitrate int) {
	if enc, ok := s.Client.AudioEncoder.(*opus.Encoder); ok {
		enc.SetBitrate(bitrate)
	}
}

func (s *State) apiAudioSetTarget(target ...*gumble.VoiceTarget) {
	if len(target) == 0 {
		s.Client.VoiceTarget = nil
		return
	}

	s.Client.Send(target[0])
	s.Client.VoiceTarget = target[0]
}

func (s *State) apiAudioIsPlaying() bool {
	s.streamMu.Lock()
	defer s.streamMu.Unlock()
	return s.stream != nil && s.stream.IsPlaying()
}

func (s *State) apiAudioCurrent() *audioStream {
	s.streamMu.Lock()
	defer s.streamMu.Unlock()
	return s.stream
}
