package plugin

import (
	"strconv"

	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
)

type UsersWrapper map[string]*gumble.User

func NewUsersWrapper(users gumble.Users) UsersWrapper {
	wrapper := UsersWrapper{}
	for _, user := range users {
		wrapper.Add(user)
	}
	return wrapper
}

func (uw UsersWrapper) Add(user *gumble.User) {
	session := strconv.FormatUint(uint64(user.Session), 10)
	uw[session] = user
}

func (uw UsersWrapper) Remove(user *gumble.User) {
	session := strconv.FormatUint(uint64(user.Session), 10)
	delete(uw, session)
}

type ChannelsWrapper map[string]*gumble.Channel

func NewChannelsWrapper(channels gumble.Channels) ChannelsWrapper {
	wrapper := ChannelsWrapper{}
	for _, channel := range channels {
		wrapper.Add(channel)
	}
	return wrapper
}

func (cw ChannelsWrapper) Add(channel *gumble.Channel) {
	id := strconv.FormatUint(uint64(channel.ID), 10)
	cw[id] = channel
}

func (cw ChannelsWrapper) Remove(channel *gumble.Channel) {
	id := strconv.FormatUint(uint64(channel.ID), 10)
	delete(cw, id)
}

type DisconnectEventWrapper struct {
	Client *gumble.Client
	Type   int

	String string

	IsError bool
	IsUser  bool

	IsOther             bool
	IsVersion           bool
	IsUserName          bool
	IsUserCredentials   bool
	IsServerPassword    bool
	IsUsernameInUse     bool
	IsServerFull        bool
	IsNoCertificate     bool
	IsAuthenticatorFail bool
}

type TextMessageEventWrapper struct {
	*gumble.TextMessageEvent
}

func (t *TextMessageEventWrapper) PlainText() string {
	return gumbleutil.PlainText(&t.TextMessageEvent.TextMessage)
}

type UserChangeEventWrapper struct {
	Client *gumble.Client
	Type   int
	User   *gumble.User
	Actor  *gumble.User

	String string

	IsConnected             bool
	IsDisconnected          bool
	IsKicked                bool
	IsBanned                bool
	IsRegistered            bool
	IsUnregistered          bool
	IsChangeName            bool
	IsChangeChannel         bool
	IsChangeComment         bool
	IsChangeAudio           bool
	IsChangeTexture         bool
	IsChangePrioritySpeaker bool
	IsChangeRecording       bool
}

type ChannelChangeEventWrapper struct {
	Client  *gumble.Client
	Type    int
	Channel *gumble.Channel

	IsCreated           bool
	IsRemoved           bool
	IsMoved             bool
	IsChangeName        bool
	IsChangeDescription bool
	IsChangePosition    bool
}

type PermissionDeniedEventWrapper struct {
	Client  *gumble.Client
	Type    int
	Channel *gumble.Channel
	User    *gumble.User

	Permission int
	String     string

	IsOther              bool
	IsPermission         bool
	IsSuperUser          bool
	IsInvalidChannelName bool
	IsTextTooLong        bool
	IsTemporaryChannel   bool
	IsMissingCertificate bool
	IsInvalidUserName    bool
	IsChannelFull        bool
	IsNestingLimit       bool
}
