package plugin

import (
	"github.com/layeh/gopher-luar"
	"github.com/layeh/gumble/gumble"
	. "github.com/layeh/piepan/plugins"
	"github.com/yuin/gopher-lua"
)

func (p *Plugin) OnConnect(e *gumble.ConnectEvent) {
	pp := p.state.GetGlobal("piepan").(*lua.LTable)
	pp.RawSetH(lua.LString("Self"), luar.New(p.state, e.Client.Self))
	pp.RawSetH(lua.LString("Users"), luar.New(p.state, e.Client.Users))
	pp.RawSetH(lua.LString("Channels"), luar.New(p.state, e.Client.Channels))

	event := ConnectEventWrapper{
		Client: e.Client,
	}
	if e.WelcomeMessage != nil {
		event.WelcomeMessage = *e.WelcomeMessage
	}
	if e.MaximumBitrate != nil {
		event.MaximumBitrate = *e.MaximumBitrate
	}

	for _, listener := range p.listeners["connect"] {
		p.callValue(listener, event)
	}
}

func (p *Plugin) OnDisconnect(e *gumble.DisconnectEvent) {
	event := DisconnectEventWrapper{
		Client: e.Client,
		Type:   int(e.Type),

		String: e.String,

		IsError: e.Type.Has(gumble.DisconnectError),
		IsUser:  e.Type.Has(gumble.DisconnectUser),

		IsOther:             e.Type.Has(gumble.DisconnectOther),
		IsVersion:           e.Type.Has(gumble.DisconnectVersion),
		IsUserName:          e.Type.Has(gumble.DisconnectUserName),
		IsUserCredentials:   e.Type.Has(gumble.DisconnectUserCredentials),
		IsServerPassword:    e.Type.Has(gumble.DisconnectServerPassword),
		IsUsernameInUse:     e.Type.Has(gumble.DisconnectUsernameInUse),
		IsServerFull:        e.Type.Has(gumble.DisconnectServerFull),
		IsNoCertificate:     e.Type.Has(gumble.DisconnectNoCertificate),
		IsAuthenticatorFail: e.Type.Has(gumble.DisconnectAuthenticatorFail),
	}

	for _, listener := range p.listeners["disconnect"] {
		p.callValue(listener, &event)
	}
}

func (p *Plugin) OnTextMessage(e *gumble.TextMessageEvent) {
	event := TextMessageEventWrapper{
		TextMessageEvent: e,
	}
	for _, listener := range p.listeners["message"] {
		p.callValue(listener, &event)
	}
}

func (p *Plugin) OnUserChange(e *gumble.UserChangeEvent) {
	event := UserChangeEventWrapper{
		Client: e.Client,
		Type:   int(e.Type),
		User:   e.User,
		Actor:  e.Actor,

		String: e.String,

		IsConnected:             e.Type.Has(gumble.UserChangeConnected),
		IsDisconnected:          e.Type.Has(gumble.UserChangeDisconnected),
		IsKicked:                e.Type.Has(gumble.UserChangeKicked),
		IsBanned:                e.Type.Has(gumble.UserChangeBanned),
		IsRegistered:            e.Type.Has(gumble.UserChangeRegistered),
		IsUnregistered:          e.Type.Has(gumble.UserChangeUnregistered),
		IsChangeName:            e.Type.Has(gumble.UserChangeName),
		IsChangeChannel:         e.Type.Has(gumble.UserChangeChannel),
		IsChangeComment:         e.Type.Has(gumble.UserChangeComment),
		IsChangeAudio:           e.Type.Has(gumble.UserChangeAudio),
		IsChangeTexture:         e.Type.Has(gumble.UserChangeTexture),
		IsChangePrioritySpeaker: e.Type.Has(gumble.UserChangePrioritySpeaker),
		IsChangeRecording:       e.Type.Has(gumble.UserChangeRecording),
	}

	for _, listener := range p.listeners["userchange"] {
		p.callValue(listener, &event)
	}
}

func (p *Plugin) OnChannelChange(e *gumble.ChannelChangeEvent) {
	event := ChannelChangeEventWrapper{
		Client:  e.Client,
		Type:    int(e.Type),
		Channel: e.Channel,

		IsCreated:           e.Type.Has(gumble.ChannelChangeCreated),
		IsRemoved:           e.Type.Has(gumble.ChannelChangeRemoved),
		IsMoved:             e.Type.Has(gumble.ChannelChangeMoved),
		IsChangeName:        e.Type.Has(gumble.ChannelChangeName),
		IsChangeDescription: e.Type.Has(gumble.ChannelChangeDescription),
		IsChangePosition:    e.Type.Has(gumble.ChannelChangePosition),
	}

	for _, listener := range p.listeners["channelchange"] {
		p.callValue(listener, &event)
	}
}

func (p *Plugin) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	event := PermissionDeniedEventWrapper{
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

	for _, listener := range p.listeners["permissiondenied"] {
		p.callValue(listener, &event)
	}
}

func (p *Plugin) OnUserList(e *gumble.UserListEvent) {
}

func (p *Plugin) OnACL(e *gumble.ACLEvent) {
}

func (p *Plugin) OnBanList(e *gumble.BanListEvent) {
}

func (p *Plugin) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
}

func (p *Plugin) OnServerConfig(e *gumble.ServerConfigEvent) {
}
