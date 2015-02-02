package plugin

import (
	"github.com/layeh/gopus"
	"github.com/layeh/gumble/gumble"
	"github.com/robertkrimen/otto"
)

func (in *Instance) apiAudioPlay(call otto.FunctionCall) otto.Value {
	if in.audio.IsPlaying() {
		return otto.FalseValue()
	}
	obj := call.Argument(0).Object()
	if obj == nil {
		return otto.FalseValue()
	}

	filenameValue, _ := obj.Get("filename")
	callbackValue, _ := obj.Get("callback")

	if callbackValue.IsFunction() {
		in.audio.Done = func() {
			in.audio.Done = nil
			in.callValue(callbackValue)
		}
	}

	if enc := in.client.AudioEncoder(); enc != nil {
		enc.SetApplication(gopus.Audio)
	}

	in.audio.Play(filenameValue.String())
	return otto.TrueValue()
}

func (in *Instance) apiAudioNewTarget(call otto.FunctionCall) otto.Value {
	id, err := call.Argument(0).ToInteger()
	if err != nil {
		return otto.UndefinedValue()
	}

	target := &gumble.VoiceTarget{}
	target.SetID(int(id))
	value, _ := in.state.ToValue(target)
	return value
}

func (in *Instance) apiAudioBitrate(call otto.FunctionCall) otto.Value {
	encoder := in.client.AudioEncoder()
	value, _ := in.state.ToValue(encoder.Bitrate())
	return value
}

func (in *Instance) apiAudioSetBitrate(call otto.FunctionCall) otto.Value {
	bitrate, err := call.Argument(0).ToInteger()
	if err != nil {
		return otto.UndefinedValue()
	}
	in.client.AudioEncoder().SetBitrate(int(bitrate))
	return otto.UndefinedValue()
}

func (in *Instance) apiAudioVolume(call otto.FunctionCall) otto.Value {
	value, _ := in.state.ToValue(in.audio.Volume)
	return value
}

func (in *Instance) apiAudioSetVolume(call otto.FunctionCall) otto.Value {
	volume, err := call.Argument(0).ToFloat()
	if err != nil {
		return otto.UndefinedValue()
	}
	in.audio.Volume = float32(volume)
	return otto.UndefinedValue()
}

func (in *Instance) apiAudioSetTarget(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		in.client.SetVoiceTarget(nil)
		return otto.TrueValue()
	}
	target, err := call.Argument(0).Export()
	if err != nil {
		return otto.UndefinedValue()
	}
	voiceTarget := target.(*gumble.VoiceTarget)
	in.client.Send(voiceTarget)
	in.client.SetVoiceTarget(voiceTarget)
	return otto.TrueValue()
}

func (in *Instance) apiAudioStop(call otto.FunctionCall) otto.Value {
	in.audio.Stop()
	return otto.UndefinedValue()
}

func (in *Instance) apiAudioIsPlaying(call otto.FunctionCall) otto.Value {
	if in.audio.IsPlaying() {
		return otto.TrueValue()
	} else {
		return otto.FalseValue()
	}
}
