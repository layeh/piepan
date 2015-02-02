package piepan

import (
	"sort"

	"github.com/layeh/bconf"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumble_ffmpeg"
)

var (
	PluginNames []string
	Plugins     map[string]*Plugin = map[string]*Plugin{}
)

type Instance struct {
	Client *gumble.Client
	FFmpeg *gumble_ffmpeg.Stream
}

type Plugin struct {
	Help string
	Init func(*Instance, *bconf.Block) error
}

func Register(name string, plugin *Plugin) {
	if plugin == nil {
		panic("piepan: plugin cannot be nil")
	}
	PluginNames = append(PluginNames, name)
	sort.Strings(PluginNames)
	Plugins[name] = plugin
}
