package plugin

import (
	"github.com/aarzilli/golua/lua"
	"github.com/layeh/gopus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/opus"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/stevedonovan/luar"
)

func (p *Plugin) apiAudioPlay(l *lua.State) int {
	if p.instance.Audio.IsPlaying() {
		l.PushBoolean(false)
		return 1
	}

	obj := luar.NewLuaObject(l, 1)
	filename := obj.Get("filename").(string)
	callback := obj.GetObject("callback")
	obj.Close()

	if enc, ok := p.instance.Client.AudioEncoder.(*opus.Encoder); ok {
		enc.SetApplication(gopus.Audio)
	}

	p.instance.Audio.Source = gumble_ffmpeg.SourceFile(filename)
	p.instance.Audio.Play()
	go func() {
		p.instance.Audio.Wait()
		if callback.Type != "nil" {
			p.callValue(callback)
		}
		callback.Close()
	}()

	return 0
}

func (p *Plugin) apiAudioNewTarget(l *lua.State) int {
	id := l.ToInteger(1)

	target := &gumble.VoiceTarget{}
	target.ID = uint32(id)

	obj := luar.NewLuaObjectFromValue(l, target)
	obj.Push()
	obj.Close()
	return 1
}

func (p *Plugin) apiAudioBitrate(l *lua.State) int {
	if enc, ok := p.instance.Client.AudioEncoder.(*opus.Encoder); ok {
		l.PushInteger(int64(enc.Bitrate()))
	} else {
		l.PushInteger(-1)
	}
	return 1
}

func (p *Plugin) apiAudioSetBitrate(l *lua.State) int {
	bitrate := l.ToInteger(1)
	if enc, ok := p.instance.Client.AudioEncoder.(*opus.Encoder); ok {
		enc.SetBitrate(bitrate)
	}
	return 0
}

func (p *Plugin) apiAudioVolume(l *lua.State) int {
	l.PushNumber(float64(p.instance.Audio.Volume))
	return 1
}

func (p *Plugin) apiAudioSetVolume(l *lua.State) int {
	volume := l.ToNumber(1)
	p.instance.Audio.Volume = float32(volume)
	return 0
}

func (p *Plugin) apiAudioSetTarget(l *lua.State) int {
	if l.GetTop() == 0 {
		p.instance.Client.VoiceTarget = nil
		return 0
	}

	voiceTarget, ok := luar.LuaToGo(l, nil, 1).(*gumble.VoiceTarget)
	if !ok {
		l.PushBoolean(false)
		return 1
	}
	p.instance.Client.Send(voiceTarget)
	p.instance.Client.VoiceTarget = voiceTarget
	return 0
}

func (p *Plugin) apiAudioStop(l *lua.State) int {
	p.instance.Audio.Stop()
	return 0
}

func (p *Plugin) apiAudioIsPlaying(l *lua.State) int {
	l.PushBoolean(p.instance.Audio.IsPlaying())
	return 1
}
