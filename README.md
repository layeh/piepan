# [piepan][1]: a bot framework for Mumble

piepan is an easy to use bot framework for interacting with a [Mumble](http://mumble.sourceforge.net/) server.  Here is a simple script that will echo back any chat message that is sent to it:

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

## Programming reference

The following section describes the API that is available for piepan scripts.

### Types

Documentation for types not part of piepan itself (e.g. User and Channel) can be found in the [gumble documentation](https://godoc.org/github.com/layeh/gumble/gumble).

#### `piepan.Audio`

- `void Play(string filename)`: Plays the given media file.
- `void SetCallback(function callback)`: Sets a function to be called after an media file is done playing. Passing nil will remove the callback function.
- `void Stop()`: Stops the currently playing stream.
- `bool IsPlaying()`: Returns true if an stream is currently playing, false otherwise.

#### `piepan.Timer`

- `piepan.Timer New(function func, int timeout)`: Creates a new timer.  After `timeout` seconds elapses, `func` will be executed.

- `void Cancel()`: Cancels the timer.

#### `piepan.Process`

- `piepan.Process New(function callback, string command, string arguments...)`: Executes `command` in a new process with the given arguments. The function `callback` is executed once the process has completed, passing if the execution was successful and the contents of standard output.

- `void Kill()`: Kills the process.

### Variables

#### [`Users`](https://godoc.org/github.com/layeh/gumble/gumble#Users) `piepan.Users`

Table containing each connected user on the server, with the keys being the session ID of the user and the value being their corresponding `piepan.User` table.

Example:

    -- prints the usernames of all the users connected to the server to standard
    -- output
    for _, user in pairs(piepan.Users) do
        print (user.Name())
    end

#### [`Channels`](https://godoc.org/github.com/layeh/gumble/gumble#Channels) `piepan.Channels`

Table that contains all of the channels that are on the server. The channels are mapped by their channel IDs. `piepan.Channels[0]` is the server's root channel.

#### [`User`](https://godoc.org/github.com/layeh/gumble/gumble#User) `piepan.Self`

The `User` table that references yourself.

### Functions

#### `piepan.Disconnect()`

Disconnects from the server.

#### `piepan.On(string event, function callback)`

Registers an event listener for a given event type. The follow events are currently supported:

- `connect` (Arguments: [`ConnectEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#ConnectEvent))
    - Called when connection to the server has been made. This is where a script should perform its initialization.
- `disconnect`
    - Called when connection to the server has been lost or after `piepan.disconnect()` is called.
- `message` (Arguments: [`TextMessageEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#TextMessageEvent))
    - Called when a text message is received.
- `userChange` (Arguments: [`UserChangeEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#UserChangeEvent))
    - Called when a user's properties changes (e.g. connects to the server).
- `channelChange` (Arguments: [`ChannelChangeEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#ChannelChangeEvent))
    - Called when a channel changes state (e.g. is added or removed).
- `permissionDenied` (Arguments: [`PermissionDenied event`](https://godoc.org/github.com/layeh/gumble/gumble#PermissionDeniedEvent))
    - Called when a requested action could not be performed.

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
