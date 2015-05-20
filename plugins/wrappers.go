package plugin

import (
	"github.com/layeh/gumble/gumble"
	"github.com/layeh/gumble/gumbleutil"
)

type ConnectEventWrapper struct {
	Client         *gumble.Client
	WelcomeMessage string
	MaximumBitrate int
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
