package piepan

import (
	"sort"
)

var (
	PluginNames []string
	Plugins     map[string]*Plugin = map[string]*Plugin{}
)

type Plugin struct {
	New func(*Instance) Environment
}

func Register(name string, plugin *Plugin) {
	if plugin == nil {
		panic("piepan: plugin cannot be nil")
	}
	PluginNames = append(PluginNames, name)
	sort.Strings(PluginNames)
	Plugins[name] = plugin
}
