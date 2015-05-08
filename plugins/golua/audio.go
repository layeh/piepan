package plugin

import (
	"github.com/layeh/gopus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/yuin/gopher-lua"
)

func (p *Plugin) apiAudioPlay(tbl *lua.LTable) bool {
	if p.instance.Audio.IsPlaying() {
		return false
	}

	filename := tbl.RawGetH(lua.LString("filename")).String()
	callback := tbl.RawGetH(lua.LString("callback"))

	if enc := p.instance.Client.AudioEncoder; enc != nil {
		enc.SetApplication(gopus.Audio)
	}

	p.instance.Audio.Source = gumble_ffmpeg.SourceFile(filename)
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
	encoder := p.instance.Client.AudioEncoder
	return encoder.Bitrate()
}

func (p *Plugin) apiAudioSetBitrate(bitrate int) {
	p.instance.Client.AudioEncoder.SetBitrate(bitrate)
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
