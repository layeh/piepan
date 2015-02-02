package plugin

import (
	"github.com/layeh/bconf"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/layeh/piepan"
)

const helpString = ` Automatically sets the audio bitrate based on the server's maximum bitrate.`

func init() {
	piepan.Register("autobitrate", &piepan.Plugin{
		Help: helpString,
		Init: func(instance *piepan.Instance, conf *bconf.Block) error {
			instance.Client.Attach(gumbleutil.AutoBitrate)
			return nil
		},
	})
}
