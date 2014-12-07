# [piepan][1]: a bot framework for Mumble

piepan is an easy to use bot framework for interacting with a [Mumble](http://mumble.sourceforge.net/) server using Lua scripts.  Here is a simple script that will echo back any chat message that is sent to it:

    -- echo.lua
    piepan.On('message', function(e)
        piepan.Self().Channel().Send(e.Message)
    end)

The above script can be started from the command line:

    $ piepan echo.lua

## Usage

    usage: piepan [options] [scripts...]
    a bot framework for Mumble
      -certificate="": user certificate file (PEM)
      -insecure=false: skip certificate checking
      -key="": user certificate key file (PEM)
      -password="": user password
      -server="localhost:64738": address of the server
      -username="piepan-bot": username of the bot

## Building

### Ubuntu 14.04

    # 1. Install dependencies
    # 1.a. Base dependencies
    sudo apt-get install -y liblua5.1-0-dev git libopus-dev mercurial wget python-software-properties software-properties-common

    # 1.b. Latest Go version (https://golang.org/dl/)
    wget https://storage.googleapis.com/golang/go1.3.3.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.3.3.linux-amd64.tar.gz
    export PATH="/usr/local/go/bin:$PATH"

    # 1.c. ffmpeg, if you would like to play media files
    sudo add-apt-repository -y ppa:jon-severinsson/ffmpeg
    sudo apt-get update
    sudo apt-get install -y ffmpeg

    # 2. Create a GOPATH (skip if you already have a GOPATH you want to use)
    export GOPATH=$(mktemp -d)

    # 3. Build piepan
    go get -u github.com/layeh/piepan/cmd/piepan

    # 4. Copy binary to current directory
    cp "$GOPATH/bin/piepan" .

## Programming reference

The following section describes the API that is available for piepan Lua scripts.

piepan is built using the [gumble](https://github.com/layeh/gumble) library. Documentation for types not part of piepan itself (including User and Channel) can be found in the [gumble documentation](https://godoc.org/github.com/layeh/gumble/gumble).

### `piepan.Audio`

- `void Play(table obj)`: Plays the media file `obj.filename`. `obj.callback` can be defined as a function that is called after the playback has completed.
- `void SetTarget(Channel|User targets...)` sets the target of subsequent `piepan.Audio.Play()` calls. Call this function with no arguments to remove any voice targeting.
- `void Stop()`: Stops the currently playing stream.
- `bool IsPlaying()`: Returns true if an stream is currently playing, false otherwise.

### [`Channels`](https://godoc.org/github.com/layeh/gumble/gumble#Channels) `piepan.Channels`

Table that contains all of the channels that are on the server. The channels are mapped by their channel IDs. `piepan.Channels[0]` is the server's root channel.

### `piepan.Disconnect()`

Disconnects from the server.

### `piepan.On(string event, function callback)`

Registers an event listener for a given event type. The follow events are currently supported:

- `connect` (Arguments: [`ConnectEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#ConnectEvent))
    - Called when connection to the server has been made. This is where a script should perform its initialization.
- `disconnect` (Arguments: [`DisconnectEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#DisconnectEvent))
    - Called when connection to the server has been lost or after `piepan.Disconnect()` is called.
- `message` (Arguments: [`TextMessageEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#TextMessageEvent))
    - Called when a text message is received.
- `userChange` (Arguments: [`UserChangeEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#UserChangeEvent))
    - Called when a user's properties changes (e.g. connects to the server).
- `channelChange` (Arguments: [`ChannelChangeEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#ChannelChangeEvent))
    - Called when a channel changes state (e.g. is added or removed).
- `permissionDenied` (Arguments: [`PermissionDenied event`](https://godoc.org/github.com/layeh/gumble/gumble#PermissionDeniedEvent))
    - Called when a requested action could not be performed.

Note: events with a `Type` field have slight changes than what is documented in gumble:

1. The `Type` field is changed to a number.
2. Individual bit flag values are added to the event as booleans prefixed with `Is`
    - Example: When a user connects to the server, `UserChangeEvent.IsConnected` will be true.

### `piepan.Process`

- `piepan.Process New(function callback, string command, string arguments...)`: Executes `command` in a new process with the given arguments. The function `callback` is executed once the process has completed, passing if the execution was successful and the contents of standard output.

- `void Kill()`: Kills the process.

### [`User`](https://godoc.org/github.com/layeh/gumble/gumble#User) `piepan.Self`

The `User` table that references yourself.

### `piepan.Timer`

- `piepan.Timer New(function callback, int timeout)`: Creates a new timer.  After `timeout` seconds elapses, `callback` will be executed.

- `void Cancel()`: Cancels the timer.

### [`Users`](https://godoc.org/github.com/layeh/gumble/gumble#Users) `piepan.Users`

Table containing each connected user on the server, with the keys being the session ID of the user and the value being their corresponding `piepan.User` table.

Example:

    -- prints the usernames of all the users connected to the server to standard
    -- output
    for _, user in pairs(piepan.Users) do
        print (user.Name())
    end

## Changelog

- Next
    - Moved to Go (+ gumble)
    - API has been overhauled. There is no backwards capability with previous versions of piepan.
- 0.3.1 (2014-10-06)
    - Fixed audio transmission memory leak
- 0.3.0 (2014-10-01)
    - Removed `data` argument from `Channel.play` and `Timer.new`
    - Fixed inability to start playing audio from inside of an audio completion callback
- 0.2.0 (2014-09-15)
    - Added support for fetching large channel descriptions, user textures/avatars, and user comments
    - Added `piepan.server.maxMessageLength`, `piepan.server.maxImageMessageLength`
    - Added `piepan.onPermissionDenied()`, `piepan.Permissions`, `piepan.PermissionDenied`, `piepan.User.hash`, `piepan.User.setComment()`, `piepan.User.register()`, `piepan.User.setTexture()`, `piepan.User.isPrioritySpeaker`, `piepan.Channel.remove()`, `piepan.Channel.setDescription()`
    - `UserChange` and `ChannelChange` are no longer hidden
    - Added audio file support
    - Each script is now loaded in its own Lua environment, preventing global variable interference between scripts
    - Fixed `piepan.User.userId` not being filled
    - Multiple instances of the same script can now be run at the same time
    - Command line option `-pw` was renamed to `-p`. The old functionality of `-p` was moved to the `-h` option, by using the new `host:port` syntax
- 0.1.1 (2013-12-15)
    - Fixed bug where event loop would stop when a packet with length zero was received
    - `sendPacket` now properly accepts a version packet
- 0.1 (2013-12-11)
    - Initial Release

## Requirements

- [gumble](https://github.com/bontibon/gumble/tree/master/gumble)
- [gumble_ffmpeg](https://github.com/bontibon/gumble/tree/master/gumble_ffmpeg)
- [golua](https://github.com/aarzilli/golua)
- [luar](https://github.com/stevedonovan/luar)

## License

This software is released under the MIT license (see LICENSE).

---

Author: Tim Cooper <<tim.cooper@layeh.com>>

[1]: https://github.com/layeh/piepan
