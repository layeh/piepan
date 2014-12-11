package piepan

import (
	"github.com/layeh/gumble/gumble"
)

func (in *Instance) OnConnect(e *gumble.ConnectEvent) {
	global, _ := in.state.Get("piepan")
	if obj := global.Object(); obj != nil {
		in.users = newUsersWrapper(in.client.Users())
		in.channels = newChannelsWrapper(in.client.Channels())

		obj.Set("Self", e.Client.Self())
		obj.Set("Users", in.users)
		obj.Set("Channels", in.channels)
	}

	for _, listener := range in.listeners["connect"] {
		in.callValue(listener, e)
	}
}

func (in *Instance) OnDisconnect(e *gumble.DisconnectEvent) {
	event := disconnectEventWrapper{
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

	in.users = nil
	in.channels = nil

	for _, listener := range in.listeners["disconnect"] {
		in.callValue(listener, &event)
	}
}

func (in *Instance) OnTextMessage(e *gumble.TextMessageEvent) {
	for _, listener := range in.listeners["message"] {
		in.callValue(listener, e)
	}
}

func (in *Instance) OnUserChange(e *gumble.UserChangeEvent) {
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

	if event.IsConnected {
		in.users.add(e.User)
	} else if event.IsDisconnected {
		in.users.remove(e.User)
	}

	for _, listener := range in.listeners["userchange"] {
		in.callValue(listener, &event)
	}
}

func (in *Instance) OnChannelChange(e *gumble.ChannelChangeEvent) {
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
		in.callValue(listener, &event)
	}
}

func (in *Instance) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
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
		in.callValue(listener, &event)
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
