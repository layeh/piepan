package piepan

import (
	"github.com/layeh/gopher-luar"
	"github.com/layeh/gumble/gumble"
)

func (s *State) OnConnect(e *gumble.ConnectEvent) {
	s.Client = e.Client

	s.table.RawSetString("Self", luar.New(s.LState, e.Client.Self))
	s.table.RawSetString("Users", luar.New(s.LState, e.Client.Users))
	s.table.RawSetString("Channels", luar.New(s.LState, e.Client.Channels))

	event := ConnectEventWrapper{
		Client: e.Client,
	}
	if e.WelcomeMessage != nil {
		event.WelcomeMessage = *e.WelcomeMessage
	}
	if e.MaximumBitrate != nil {
		event.MaximumBitrate = *e.MaximumBitrate
	}

	for _, listener := range s.listeners["connect"] {
		s.callValue(listener, event)
	}
}

func (s *State) OnDisconnect(e *gumble.DisconnectEvent) {
	event := DisconnectEventWrapper{
		Client: e.Client,
		Type:   int(e.Type),

		String: e.String,

		IsError: e.Type.Has(gumble.DisconnectError),
		IsUser:  e.Type.Has(gumble.DisconnectUser),
	}

	for _, listener := range s.listeners["disconnect"] {
		s.callValue(listener, &event)
	}
}

func (s *State) OnTextMessage(e *gumble.TextMessageEvent) {
	event := TextMessageEventWrapper{
		TextMessageEvent: e,
	}
	for _, listener := range s.listeners["message"] {
		s.callValue(listener, &event)
	}
}

func (s *State) OnUserChange(e *gumble.UserChangeEvent) {
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

	for _, listener := range s.listeners["userchange"] {
		s.callValue(listener, &event)
	}
}

func (s *State) OnChannelChange(e *gumble.ChannelChangeEvent) {
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

	for _, listener := range s.listeners["channelchange"] {
		s.callValue(listener, &event)
	}
}

func (s *State) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
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

	for _, listener := range s.listeners["permissiondenied"] {
		s.callValue(listener, &event)
	}
}

func (s *State) OnUserList(e *gumble.UserListEvent) {
}

func (s *State) OnACL(e *gumble.ACLEvent) {
}

func (s *State) OnBanList(e *gumble.BanListEvent) {
}

func (s *State) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
}

func (s *State) OnServerConfig(e *gumble.ServerConfigEvent) {
}
