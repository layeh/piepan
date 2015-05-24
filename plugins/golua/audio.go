package plugin

import (
	"github.com/layeh/gopus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/opus"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/yuin/gopher-lua"
)

func (p *Plugin) apiAudioPlay(tbl *lua.LTable) bool {
	if p.instance.Audio.IsPlaying() {
		return false
	}

	filename := tbl.RawGetH(lua.LString("filename"))

	exec := tbl.RawGetH(lua.LString("exec"))
	args := tbl.RawGetH(lua.LString("args"))

	callback := tbl.RawGetH(lua.LString("callback"))

	if enc, ok := p.instance.Client.AudioEncoder.(*opus.Encoder); ok {
		enc.SetApplication(gopus.Audio)
	}

	switch {
	// source file
	case filename != lua.LNil && exec == lua.LNil && args == lua.LNil:
		p.instance.Audio.Source = gumble_ffmpeg.SourceFile(filename.String())
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
		p.instance.Audio.Source = gumble_ffmpeg.SourceExec(exec.String(), argsStr...)
	default:
		panic("invalid source type")
	}

	p.instance.Audio.Play()
	go func() {
		p.instance.Audio.Wait()
		if callback.Type() != lua.LTNil {
			p.callValue(callback)
		}
	}()

	return true
}

func (p *Plugin) apiAudioNewTarget(id uint32) *gumble.VoiceTarget {
	return &gumble.VoiceTarget{
		ID: id,
	}
}

func (p *Plugin) apiAudioBitrate() int {
	if enc, ok := p.instance.Client.AudioEncoder.(*opus.Encoder); ok {
		return enc.Bitrate()
	}
	return -1
}

func (p *Plugin) apiAudioSetBitrate(bitrate int) {
	if enc, ok := p.instance.Client.AudioEncoder.(*opus.Encoder); ok {
		enc.SetBitrate(bitrate)
	}
}

func (p *Plugin) apiAudioVolume() float32 {
	return p.instance.Audio.Volume
}

func (p *Plugin) apiAudioSetVolume(volume float32) {
	p.instance.Audio.Volume = volume
}

func (p *Plugin) apiAudioSetTarget(target ...*gumble.VoiceTarget) {
	if len(target) == 0 {
		p.instance.Client.VoiceTarget = nil
		return
	}

	p.instance.Client.Send(target[0])
	p.instance.Client.VoiceTarget = target[0]
}

func (p *Plugin) apiAudioStop() {
	p.instance.Audio.Stop()
}

func (p *Plugin) apiAudioIsPlaying() bool {
	return p.instance.Audio.IsPlaying()
}
