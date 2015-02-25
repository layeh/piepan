package piepan

import (
	"sort"
)

type PluginExts []string

func (p PluginExts) Len() int {
	return len(p)
}

func (p PluginExts) Less(i, j int) bool {
	if len(p[i]) == len(p[j]) {
		return p[i] < p[j]
	}
	return len(p[i]) > len(p[j])
}

func (p PluginExts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

var (
	PluginExtensions PluginExts
	Plugins          map[string]*Plugin = map[string]*Plugin{}
)

type Plugin struct {
	Name string
	New  func(*Instance) Environment
}

func Register(extension string, plugin *Plugin) {
	if plugin == nil {
		panic("piepan: plugin cannot be nil")
	}
	if _, ok := Plugins[extension]; ok {
		panic("piepan: extension " + extension + " already registered")
	}
	PluginExtensions = append(PluginExtensions, extension)
	sort.Sort(PluginExtensions)
	Plugins[extension] = plugin
}
