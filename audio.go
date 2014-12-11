package piepan

import (
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

	in.audio.Play(filenameValue.String())
	return otto.TrueValue()
}

func (in *Instance) apiAudioSetTarget(call otto.FunctionCall) otto.Value {
	if len(call.ArgumentList) == 0 {
		in.client.SetVoiceTarget(nil)
		return otto.UndefinedValue()
	}

	vt := gumble.VoiceTarget{}
	vt.SetID(1)
	for _, arg := range call.ArgumentList {
		value, _ := arg.Export()
		switch val := value.(type) {
		case *gumble.User:
			vt.AddUser(val)
		case *gumble.Channel:
			vt.AddChannel(val, false, false)
		}
	}
	in.client.Send(&vt)
	in.client.SetVoiceTarget(&vt)

	return otto.UndefinedValue()
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
