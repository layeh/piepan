package piepan

import (
	"github.com/layeh/gopus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleffmpeg"
	"github.com/layeh/gumble/opus"
	"github.com/yuin/gopher-lua"
)

func (s *State) apiAudioPlay(tbl *lua.LTable) bool {
	filename := tbl.RawGetString("filename")

	exec := tbl.RawGetString("exec")
	args := tbl.RawGetString("args")

	callback := tbl.RawGetString("callback")

	if s.audioStream != nil {
		s.audioStream.Stop()
	}

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
		panic("invalid piepan.Audio.Play source type")
	}

	s.audioStream = gumbleffmpeg.New(s.Client, source)
	s.audioStream.Play()
	go func(stream *gumbleffmpeg.Stream) {
		stream.Wait()
		if callback.Type() != lua.LTNil {
			s.callValue(callback)
		}
	}(s.audioStream)

	return true
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

func (s *State) apiAudioVolume() float32 {
	if s.audioStream != nil {
		return s.audioStream.Volume
	}
	return s.audioVolume
}

func (s *State) apiAudioSetVolume(volume float32) {
	s.audioVolume = volume
	if s.audioStream != nil {
		s.audioStream.Volume = volume
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

func (s *State) apiAudioStop() {
	if s.audioStream != nil {
		s.audioStream.Stop()
	}
}

func (s *State) apiAudioIsPlaying() bool {
	if s.audioStream == nil {
		return false
	}
	return s.audioStream.State() == gumbleffmpeg.StatePlaying
}
