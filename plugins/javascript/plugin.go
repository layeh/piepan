package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/layeh/piepan"
	. "github.com/layeh/piepan/plugins"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
)

func init() {
	piepan.Register("js", &piepan.Plugin{
		Name: "JavaScript",
		New: func(in *piepan.Instance) piepan.Environment {
			s := otto.New()
			p := &Plugin{
				instance:  in,
				state:     s,
				listeners: make(map[string][]otto.Value),
			}
			s.Set("piepan", map[string]interface{}{
				"On":         p.apiOn,
				"Disconnect": p.apiDisconnect,
				"Audio": map[string]interface{}{
					"Play":       p.apiAudioPlay,
					"IsPlaying":  p.apiAudioIsPlaying,
					"Stop":       p.apiAudioStop,
					"NewTarget":  p.apiAudioNewTarget,
					"SetTarget":  p.apiAudioSetTarget,
					"Bitrate":    p.apiAudioBitrate,
					"SetBitrate": p.apiAudioSetBitrate,
					"Volume":     p.apiAudioVolume,
					"SetVolume":  p.apiAudioSetVolume,
				},
				"File": map[string]interface{}{
					"Open": p.apiFileOpen,
				},
				"Process": map[string]interface{}{
					"New": p.apiProcessNew,
				},
				"Timer": map[string]interface{}{
					"New": p.apiTimerNew,
				},
			})
			s.Set("ENV", createEnvVars())
			return p
		},
	})
}

type Plugin struct {
	instance *piepan.Instance

	stateLock sync.Mutex
	state     *otto.Otto

	listeners map[string][]otto.Value
	users     UsersWrapper
	channels  ChannelsWrapper
}

func (p *Plugin) LoadScriptFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	script, err := p.state.Compile(filename, data)
	if err != nil {
		return err
	}
	_, err = p.state.Run(script)
	if err != nil {
		return err
	}
	return nil
}

func (p *Plugin) apiOn(call otto.FunctionCall) otto.Value {
	event := strings.ToLower(call.Argument(0).String())
	function := call.Argument(1)
	p.listeners[event] = append(p.listeners[event], function)
	return otto.UndefinedValue()
}

func (p *Plugin) apiDisconnect() {
	if client := p.instance.Client; client != nil {
		client.Disconnect()
	}
}

func (p *Plugin) error(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
}

func (p *Plugin) callValue(value otto.Value, arguments ...interface{}) {
	for i, arg := range arguments {
		value, err := p.state.ToValue(arg)
		if err != nil {
			p.error(err)
			return
		}
		arguments[i] = value
	}
	p.stateLock.Lock()
	if _, err := value.Call(otto.NullValue(), arguments...); err != nil {
		p.stateLock.Unlock()
		p.error(err)
	} else {
		p.stateLock.Unlock()
	}
}
