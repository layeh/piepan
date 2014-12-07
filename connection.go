package piepan

import (
	"fmt"
	"os"

	"github.com/layeh/gumble/gumble"
	"github.com/stevedonovan/luar"
)

func (in *Instance) OnConnect(e *gumble.ConnectEvent) {
	luar.Register(in.state, "piepan", luar.Map{
		"Self":       e.Client.Self(),
		"Users":      e.Client.Users(),
		"Channels":   e.Client.Channels(),
		"Disconnect": in.disconnect,
	})

	in.stateLock.Lock()
	defer in.stateLock.Unlock()

	for _, listener := range in.listeners["connect"] {
		if _, err := listener.Call(e); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func (in *Instance) OnDisconnect(e *gumble.DisconnectEvent) {
	in.stateLock.Lock()
	defer in.stateLock.Unlock()

	for _, listener := range in.listeners["disconnect"] {
		if _, err := listener.Call(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func (in *Instance) OnTextMessage(e *gumble.TextMessageEvent) {
	in.stateLock.Lock()
	defer in.stateLock.Unlock()

	for _, listener := range in.listeners["message"] {
		if _, err := listener.Call(e); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func (in *Instance) OnUserChange(e *gumble.UserChangeEvent) {
	in.stateLock.Lock()
	defer in.stateLock.Unlock()

	for _, listener := range in.listeners["userchange"] {
		if _, err := listener.Call(e); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func (in *Instance) OnChannelChange(e *gumble.ChannelChangeEvent) {
	in.stateLock.Lock()
	defer in.stateLock.Unlock()

	for _, listener := range in.listeners["channelchange"] {
		if _, err := listener.Call(e); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func (in *Instance) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	in.stateLock.Lock()
	defer in.stateLock.Unlock()

	for _, listener := range in.listeners["permissiondenied"] {
		if _, err := listener.Call(e); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func (in *Instance) OnUserList(e *gumble.UserListEvent) {
}

func (in *Instance) OnAcl(e *gumble.AclEvent) {
}

func (in *Instance) OnBanList(e *gumble.BanListEvent) {
}

func (in *Instance) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
}
