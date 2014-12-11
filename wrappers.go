package piepan

import (
	"strconv"

	"github.com/layeh/gumble/gumble"
)

type usersWrapper map[string]*gumble.User

func newUsersWrapper(users gumble.Users) usersWrapper {
	wrapper := usersWrapper{}
	for _, user := range users {
		wrapper.add(user)
	}
	return wrapper
}

func (uw usersWrapper) add(user *gumble.User) {
	session := strconv.FormatUint(uint64(user.Session()), 10)
	uw[session] = user
}

func (uw usersWrapper) remove(user *gumble.User) {
	session := strconv.FormatUint(uint64(user.Session()), 10)
	delete(uw, session)
}

type channelsWrapper map[string]*gumble.Channel

func newChannelsWrapper(channels gumble.Channels) channelsWrapper {
	wrapper := channelsWrapper{}
	for _, channel := range channels {
		wrapper.add(channel)
	}
	return wrapper
}

func (cw channelsWrapper) add(channel *gumble.Channel) {
	id := strconv.FormatUint(uint64(channel.ID()), 10)
	cw[id] = channel
}

func (cw channelsWrapper) remove(channel *gumble.Channel) {
	id := strconv.FormatUint(uint64(channel.ID()), 10)
	delete(cw, id)
}

type disconnectEventWrapper struct {
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

type userChangeEventWrapper struct {
	Client *gumble.Client
	Type   int
	User   *gumble.User
	Actor  *gumble.User

	String string

	IsConnected     bool
	IsDisconnected  bool
	IsKicked        bool
	IsBanned        bool
	IsChangeName    bool
	IsChangeChannel bool
	IsChangeComment bool
}

type channelChangeEventWrapper struct {
	Client  *gumble.Client
	Type    int
	Channel *gumble.Channel

	IsCreated           bool
	IsRemoved           bool
	IsMoved             bool
	IsChangeName        bool
	IsChangeDescription bool
}

type permissionDeniedEventWrapper struct {
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
