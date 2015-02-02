package plugin

import (
	"github.com/layeh/bconf"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/layeh/piepan"
)

const helpString = ` Automatically sets the audio bitrate based on the server's maximum bitrate.`

func init() {
	piepan.Register("autobitrate", &piepan.Plugin{
		Help: helpString,
		Init: func(client *gumble.Client, conf *bconf.Block) error {
			client.Attach(gumbleutil.AutoBitrate)
			return nil
		},
	})
}
