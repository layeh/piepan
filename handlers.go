package piepan

import (
	"github.com/layeh/gumble/gumble"
)

func (in *Instance) OnConnect(e *gumble.ConnectEvent) {
	for _, env := range in.envs {
		env.OnConnect(e)
	}
}

func (in *Instance) OnDisconnect(e *gumble.DisconnectEvent) {
	for _, env := range in.envs {
		env.OnDisconnect(e)
	}
}

func (in *Instance) OnTextMessage(e *gumble.TextMessageEvent) {
	for _, env := range in.envs {
		env.OnTextMessage(e)
	}
}

func (in *Instance) OnUserChange(e *gumble.UserChangeEvent) {
	for _, env := range in.envs {
		env.OnUserChange(e)
	}
}

func (in *Instance) OnChannelChange(e *gumble.ChannelChangeEvent) {
	for _, env := range in.envs {
		env.OnChannelChange(e)
	}
}
func (in *Instance) OnPermissionDenied(e *gumble.PermissionDeniedEvent) {
	for _, env := range in.envs {
		env.OnPermissionDenied(e)
	}
}
func (in *Instance) OnUserList(e *gumble.UserListEvent) {
	for _, env := range in.envs {
		env.OnUserList(e)
	}
}
func (in *Instance) OnACL(e *gumble.ACLEvent) {
	for _, env := range in.envs {
		env.OnACL(e)
	}
}
func (in *Instance) OnBanList(e *gumble.BanListEvent) {
	for _, env := range in.envs {
		env.OnBanList(e)
	}
}
func (in *Instance) OnContextActionChange(e *gumble.ContextActionChangeEvent) {
	for _, env := range in.envs {
		env.OnContextActionChange(e)
	}
}
