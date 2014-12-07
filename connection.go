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

	event := userChangeEventWrapper{
		Client: e.Client,
		Type:   int(e.Type),
		User:   e.User,
		Actor:  e.Actor,

		String: e.String,

		IsConnected:     e.Type.Has(gumble.UserChangeConnected),
		IsDisconnected:  e.Type.Has(gumble.UserChangeDisconnected),
		IsKicked:        e.Type.Has(gumble.UserChangeKicked),
		IsBanned:        e.Type.Has(gumble.UserChangeBanned),
		IsChangeName:    e.Type.Has(gumble.UserChangeName),
		IsChangeChannel: e.Type.Has(gumble.UserChangeChannel),
		IsChangeComment: e.Type.Has(gumble.UserChangeComment),
	}

	for _, listener := range in.listeners["userchange"] {
		if _, err := listener.Call(&event); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func (in *Instance) OnChannelChange(e *gumble.ChannelChangeEvent) {
	in.stateLock.Lock()
	defer in.stateLock.Unlock()

	event := channelChangeEventWrapper{
		Client:  e.Client,
		Type:    int(e.Type),
		Channel: e.Channel,

		IsCreated:           e.Type.Has(gumble.ChannelChangeCreated),
		IsRemoved:           e.Type.Has(gumble.ChannelChangeRemoved),
		IsMoved:             e.Type.Has(gumble.ChannelChangeMoved),
		IsChangeName:        e.Type.Has(gumble.ChannelChangeName),
		IsChangeDescription: e.Type.Has(gumble.ChannelChangeDescription),
	}

	for _, listener := range in.listeners["channelchange"] {
		if _, err := listener.Call(&event); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func (in *Instance) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	in.stateLock.Lock()
	defer in.stateLock.Unlock()

	event := permissionDeniedEventWrapper{
		Client:  e.Client,
		Type:    int(e.Type),
		Channel: e.Channel,
		User:    e.User,

		Permission: int(e.Permission),
		String:     e.String,

		IsOther:              e.Type.Has(gumble.PermissionDeniedOther),
		IsPermission:         e.Type.Has(gumble.PermissionDeniedPermission),
		IsSuperUser:          e.Type.Has(gumble.PermissionDeniedSuperUser),
		IsInvalidChannelName: e.Type.Has(gumble.PermissionDeniedInvalidChannelName),
		IsTextTooLong:        e.Type.Has(gumble.PermissionDeniedTextTooLong),
		IsTemporaryChannel:   e.Type.Has(gumble.PermissionDeniedTemporaryChannel),
		IsMissingCertificate: e.Type.Has(gumble.PermissionDeniedMissingCertificate),
		IsInvalidUserName:    e.Type.Has(gumble.PermissionDeniedInvalidUserName),
		IsChannelFull:        e.Type.Has(gumble.PermissionDeniedChannelFull),
		IsNestingLimit:       e.Type.Has(gumble.PermissionDeniedNestingLimit),
	}

	for _, listener := range in.listeners["permissiondenied"] {
		if _, err := listener.Call(&event); err != nil {
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
