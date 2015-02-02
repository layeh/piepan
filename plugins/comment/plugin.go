package plugin

import (
	"github.com/layeh/bconf"
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
	"github.com/layeh/piepan"
)

const helpString = ` Sets the comment for the bot upon joining a server.
  Configuration:
    comment <string>: the comment.`

func init() {
	piepan.Register("comment", &piepan.Plugin{
		Help: helpString,
		Init: func(instance *piepan.Instance, conf *bconf.Block) error {
			instance.Client.Attach(gumbleutil.Listener{
				Connect: func(e *gumble.ConnectEvent) {
					e.Client.Self().SetComment(conf.String("comment"))
				},
			})
			return nil
		},
	})
}
