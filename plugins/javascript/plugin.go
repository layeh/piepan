package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/layeh/bconf"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
	"github.com/layeh/piepan"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

const helpString = ` Scripting via JavaScript.
 Configuration:
   file <string>: the file names of scripts that will be executed. Can appear
                  multiple times in the same plugin block.`

func init() {
	piepan.Register("javascript", &piepan.Plugin{
		Help: helpString,
		Init: func(client *gumble.Client, conf *bconf.Block) error {
			instance := New(client)
			for _, script := range conf.Fields["file"] {
				if err := instance.LoadScriptFile(script.String(0)); err != nil {
					return err
				}
			}
			instance.ErrFunc = func(err error) {
				if ottoErr, ok := err.(*otto.Error); ok {
					fmt.Fprintf(os.Stderr, "%s\n", ottoErr.String())
				} else {
					fmt.Fprintf(os.Stderr, "%s\n", err)
				}
			}
			client.Attach(instance)
			return nil
		},
	})
}

type Instance struct {
	ErrFunc func(error)

	client *gumble.Client

	audio *gumble_ffmpeg.Stream

	stateLock sync.Mutex
	state     *otto.Otto
	listeners map[string][]otto.Value
	users     usersWrapper
	channels  channelsWrapper
}

func New(client *gumble.Client) *Instance {
	in := &Instance{
		client:    client,
		state:     otto.New(),
		listeners: make(map[string][]otto.Value),
	}
	in.audio, _ = gumble_ffmpeg.New(in.client)

	in.state.Set("piepan", map[string]interface{}{
		"On":         in.apiOn,
		"Disconnect": in.apiDisconnect,
		"Audio": map[string]interface{}{
			"Play":       in.apiAudioPlay,
			"IsPlaying":  in.apiAudioIsPlaying,
			"Stop":       in.apiAudioStop,
			"NewTarget":  in.apiAudioNewTarget,
			"SetTarget":  in.apiAudioSetTarget,
			"Bitrate":    in.apiAudioBitrate,
			"SetBitrate": in.apiAudioSetBitrate,
			"Volume":     in.apiAudioVolume,
			"SetVolume":  in.apiAudioSetVolume,
		},
		"Process": map[string]interface{}{
			"New": in.apiProcessNew,
		},
		"Timer": map[string]interface{}{
			"New": in.apiTimerNew,
		},
	})
	in.state.Set("ENV", in.createEnvVars())
	return in
}

func (in *Instance) createEnvVars() map[string]string {
	vars := make(map[string]string)
	for _, val := range os.Environ() {
		split := strings.SplitN(val, "=", 2)
		if len(split) != 2 {
			continue
		}
		vars[split[0]] = split[1]
	}
	return vars
}

func (in *Instance) callValue(value otto.Value, arguments ...interface{}) {
	for i, arg := range arguments {
		value, err := in.state.ToValue(arg)
		if err != nil {
			if errFunc := in.ErrFunc; errFunc != nil {
				errFunc(err)
			}
			return
		}
		arguments[i] = value
	}
	in.stateLock.Lock()
	if _, err := value.Call(otto.NullValue(), arguments...); err != nil {
		in.stateLock.Unlock()
		if errFunc := in.ErrFunc; errFunc != nil {
			errFunc(err)
		}
	} else {
		in.stateLock.Unlock()
	}
}

func (in *Instance) apiOn(call otto.FunctionCall) otto.Value {
	event := strings.ToLower(call.Argument(0).String())
	function := call.Argument(1)
	in.listeners[event] = append(in.listeners[event], function)
	return otto.UndefinedValue()
}

func (in *Instance) apiDisconnect() {
	if client := in.client; client != nil {
		client.Disconnect()
	}
}

func (in *Instance) LoadScriptFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	script, err := in.state.Compile(filename, data)
	if err != nil {
		return err
	}
	_, err = in.state.Run(script)
	if err != nil {
		return err
	}
	return nil
}
