package plugin

import (
	"github.com/layeh/gopus"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/opus"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/robertkrimen/otto"
)

func (p *Plugin) apiAudioPlay(call otto.FunctionCall) otto.Value {
	if p.instance.Audio.IsPlaying() {
		return otto.FalseValue()
	}
	obj := call.Argument(0).Object()
	if obj == nil {
		return otto.FalseValue()
	}

	filenameValue, _ := obj.Get("filename")
	callbackValue, _ := obj.Get("callback")

	if enc, ok := p.instance.Client.AudioEncoder.(*opus.Encoder); ok {
		enc.SetApplication(gopus.Audio)
	}

	p.instance.Audio.Source = gumble_ffmpeg.SourceFile(filenameValue.String())
	p.instance.Audio.Play()
	go func() {
		p.instance.Audio.Wait()
		if callbackValue.IsFunction() {
			p.callValue(callbackValue)
		}
	}()
	return otto.TrueValue()
}

func (p *Plugin) apiAudioNewTarget(call otto.FunctionCall) otto.Value {
	id, err := call.Argument(0).ToInteger()
	if err != nil {
		return otto.UndefinedValue()
	}

	target := &gumble.VoiceTarget{}
	target.ID = uint32(id)
	value, _ := p.state.ToValue(target)
	return value
}

func (p *Plugin) apiAudioBitrate(call otto.FunctionCall) otto.Value {
	if enc, ok := p.instance.Client.AudioEncoder.(*opus.Encoder); ok {
		value, _ := p.state.ToValue(enc.Bitrate())
		return value
	}
	return otto.UndefinedValue()
}

func (p *Plugin) apiAudioSetBitrate(call otto.FunctionCall) otto.Value {
	bitrate, err := call.Argument(0).ToInteger()
	if err != nil {
		return otto.UndefinedValue()
	}
	if enc, ok := p.instance.Client.AudioEncoder.(*opus.Encoder); ok {
		enc.SetBitrate(int(bitrate))
	}
	return otto.UndefinedValue()
}

func (p *Plugin) apiAudioVolume(call otto.FunctionCall) otto.Value {
	value, _ := p.state.ToValue(p.instance.Audio.Volume)
	return value
}

func (p *Plugin) apiAudioSetVolume(call otto.FunctionCall) otto.Value {
	volume, err := call.Argument(0).ToFloat()
	if err != nil {
		return otto.UndefinedValue()
	}
	p.instance.Audio.Volume = float32(volume)
	return otto.UndefinedValue()
}

func (p *Plugin) apiAudioSetTarget(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		p.instance.Client.VoiceTarget = nil
		return otto.TrueValue()
	}
	target, err := call.Argument(0).Export()
	if err != nil {
		return otto.UndefinedValue()
	}
	voiceTarget := target.(*gumble.VoiceTarget)
	p.instance.Client.Send(voiceTarget)
	p.instance.Client.VoiceTarget = voiceTarget
	return otto.TrueValue()
}

func (p *Plugin) apiAudioStop(call otto.FunctionCall) otto.Value {
	p.instance.Audio.Stop()
	return otto.UndefinedValue()
}

func (p *Plugin) apiAudioIsPlaying(call otto.FunctionCall) otto.Value {
	if p.instance.Audio.IsPlaying() {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}
