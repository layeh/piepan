# piepan: an easy to use framework for writing [Mumble](http://mumble.info) bots using [Lua](http://lua.org/)

## Usage

    piepan v0.10.0
    usage: piepan [options] [script files]
    an easy to use framework for writing Mumble bots using Lua
      -access-token value
            server access token (can be defined multiple times) (default [])
      -certificate string
            user certificate file (PEM)
      -ffmpeg string
            ffmpeg-capable executable for media streaming (default "ffmpeg")
      -insecure
            skip certificate checking
      -key string
            user certificate key file (PEM)
      -password string
            user password
      -server string
            address of the server (default "localhost:64738")
      -username string
            username of the bot (default "piepan-bot")

## Building

### Windows

1. Install dependencies
    1. Base dependencies
        1. Download and run [Go installer](https://golang.org/dl/)
        2. Download and run [MSYS2 installer](https://msys2.github.io/)
            - Uncheck "Run MSYS2 32/64bit now"
    2. Open the MSYS2 "MinGW-w64 Win32/64 Shell" from the start menu to install additional dependencies
        - 32-bit:
            - `pacman -Syy mingw-w64-i686-toolchain git mingw-w64-i686-opus pkg-config mingw-w64-i686-ffmpeg`
        - 64-bit:
            - `pacman -Syy mingw-w64-x86_64-toolchain git mingw-w64-x86_64-opus pkg-config mingw-w64-x86_64-ffmpeg`
2. Create a GOPATH (skip if you already have a GOPATH you want to use)
    - `export GOPATH=$(mktemp -d)`
3. Configure environment for building Opus
    - `export CGO_LDFLAGS="$(pkg-config --libs opus)"`
    - `export CGO_CFLAGS="$(pkg-config --cflags opus)"`
4. Fetch piepan
    - `go get -tags nopkgconfig -u layeh.com/piepan`
5. Build piepan
    - `go build -o piepan.exe $GOPATH/src/layeh.com/piepan/cmd/piepan/main.go`
6. Run piepan
    - `./piepan.exe -ffmpeg=ffmpeg.exe ...`

### Ubuntu 14.04

1. Install dependencies
  1. Base dependencies
      - `sudo apt-get install -y git libopus-dev wget python-software-properties software-properties-common pkg-config gcc libav-tools`
  2. Go
      - `sudo add-apt-repository -y ppa:evarlast/golang1.4`
      - `sudo apt-get update`
      - `sudo apt-get install -y golang`
2. Create a GOPATH (skip if you already have a GOPATH you want to use)
    - `export GOPATH=$(mktemp -d)`
3. Fetch piepan
    - `go get -u layeh.com/piepan`
4. Build piepan
    - `go build -o piepan $GOPATH/src/layeh.com/piepan/cmd/piepan/main.go`
5. Run piepan with `avconv`
    - `./piepan -ffmpeg=avconv ...`

## API

piepan is built using the [gumble](https://layeh.com/gumble) library. Documentation for types not part of piepan itself (including User and Channel) can be found in the [gumble documentation](https://godoc.org/layeh.com/gumble/gumble).

### `piepan.Audio`

- [`VoiceTarget`](https://godoc.org/layeh.com/gumble/gumble#VoiceTarget)   `NewTarget(int id)`: Create a new voice target object.
- `piepan.AudioStream New(table obj)`: Returns a new AudioStream. The following modes are supported:
    - Filename:
        - Plays the media file `obj.filename`.
    - Pipe:
        - Uses the output of the program `obj.exec` executed with `obj.args`.

    `obj.offset` defines the number of seconds from the beginning of the stream to starting playing at.  
    `obj.callback` can be defined as a function that is called after the playback has completed.
    If audio has been paused
- `piepan.AudioStream Current()`: Returns the currently playing stream.
- `void SetTarget(VoiceTarget target)` sets the target for the new audio stream that is played. Calling this function with no arguments removes any voice targeting.
- `int Bitrate()`: Returns the bitrate of the audio encoder.
- `void SetBitrate(int bitrate)`: Sets the bitrate of the audio encoder. Calling this function will override the automatically-configured, optimal bitrate.
- `bool IsPlaying()`: Returns if there is a stream currently playing.

### `piepan.AudioStream`

Note: after an audio stream has stopped/completed, it cannot be started again.

- `void Play()`: Starts playing or resumes the audio stream. An error is thrown if the stream cannot be started, usually due to a missing ffmpeg executable.
- `void Stop()`: Stops a playing or paused stream.
- `void Pause()`: Pauses the playing stream.
- `void IsPlaying()`: Returns if the stream is playing.
- `void IsPaused()`: Returns if the stream is paused.
- `void IsStopped()`: Returns if the stream has stopped.
- `float Elapsed()`: Returns the amount of audio (in seconds) that the stream played.
- `void SetVolume(float volume)`: Sets the volume of transmitted audio (default: 1.0).
- `float Volume()`: Returns the stream's volume.

### [`Channels`](https://godoc.org/layeh.com/gumble/gumble#Channels) `piepan.Channels`

Object that contains all of the channels that are on the server. The channels are mapped by their channel IDs. `piepan.Channels[0]` is the server's root channel.

### `piepan.Disconnect()`

Disconnects from the server.

### `piepan.On(string event, function callback)`

Registers an event listener for a given event type. The follow events are currently supported:

- `connect` (Arguments: [`ConnectEvent event`](https://godoc.org/layeh.com/gumble/gumble#ConnectEvent))
    - Called when connection to the server has been made. This is where a script should perform its initialization.
- `disconnect` (Arguments: [`DisconnectEvent event`](https://godoc.org/layeh.com/gumble/gumble#DisconnectEvent))
    - Called when connection to the server has been lost or after `piepan.Disconnect()` is called.
- `message` (Arguments: [`TextMessageEvent event`](https://godoc.org/layeh.com/gumble/gumble#TextMessageEvent))
    - Called when a text message is received.
- `userChange` (Arguments: [`UserChangeEvent event`](https://godoc.org/layeh.com/gumble/gumble#UserChangeEvent))
    - Called when a user's properties changes (e.g. connects to the server).
- `channelChange` (Arguments: [`ChannelChangeEvent event`](https://godoc.org/layeh.com/gumble/gumble#ChannelChangeEvent))
    - Called when a channel changes state (e.g. is added or removed).
- `permissionDenied` (Arguments: [`PermissionDeniedEvent event`](https://godoc.org/layeh.com/gumble/gumble#PermissionDeniedEvent))
    - Called when a requested action could not be performed.
- `stream` (Arguments: `piepan.AudioStream`)
    - Called when a stream changes state (plays, pauses, or stops).

Note: events with a `Type` field have slight changes than what is documented in gumble:

1. The `Type` field is changed to a number.
2. Individual bit flag values are added to the event as booleans prefixed with `Is`
    - `DisconnectEvent`
        - `IsError`
        - `IsUser`
    - `UserChangeEvent`
        - `IsConnected`
        - `IsDisconnected`
        - `IsKicked`
        - `IsBanned`
        - `IsRegistered`
        - `IsUnregistered`
        - `IsChangeName`
        - `IsChangeChannel`
        - `IsChangeComment`
        - `IsChangeAudio`
        - `IsChangeTexture`
        - `IsChangePrioritySpeaker`
        - `IsChangeRecording`
    - `ChannelChangeEvent`
        - `IsCreated`
        - `IsRemoved`
        - `IsMoved`
        - `IsChangeName`
        - `IsChangeDescription`
        - `IsChangePosition`
    - `PermissionDeniedEvent`
        - `IsOther`
        - `IsPermission`
        - `IsSuperUser`
        - `IsInvalidChannelName`
        - `IsTextTooLong`
        - `IsTemporaryChannel`
        - `IsMissingCertificate`
        - `IsInvalidUserName`
        - `IsChannelFull`
        - `IsNestingLimit`

### `piepan.Process`

- `piepan.Process New(function callback, string command, string arguments...)`: Executes `command` in a new process with the given arguments. The function `callback` is executed once the process has completed, passing if the execution was successful and the contents of standard output.

- `void Kill()`: Kills the process.

### [`User`](https://godoc.org/layeh.com/gumble/gumble#User) `piepan.Self`

The `User` object that references yourself.

### `piepan.Timer`

- `piepan.Timer New(function callback, int timeout)`: Creates a new timer.  After at least `timeout` milliseconds, `callback` will be executed.

- `void Cancel()`: Cancels the timer.

### [`Users`](https://godoc.org/layeh.com/gumble/gumble#Users) `piepan.Users`

Object containing each connected user on the server, with the keys being the session ID of the user and the value being their corresponding `piepan.User` table.

### `piepan.Args`

- `piepan.Args`: features arguments passed to piepan via the `-script-args` flag as an array. 0 or more `-script-args` may be included

## Changelog

- 0.10.0 (Next)
    - Remove -lock flag
    - Server connection rejection is now earlier in the lifecycle
    - Change exit codes
    - Added `-script-args` which surfaces in lua as `piepan.Args`
- 0.9.0 (2016-01-15)
    - Add "stream" event.
    - Fix `piepan.Timer.New` not being exposed
- 0.8.1 (2015-12-16)
    - Fix -ffmpeg flag not being used
    - `AudioStream.Play` now throws an error if the stream cannot be started
- 0.8.0 (2015-12-11)
    - Switch to stream-based audio. Individual audios streams can be created then played, paused, and stopped.
    - Add `gumble.ConnectEvent` wrapper
    - Add pipe support to `piepan.Audio.New`
    - Add `offset` field to `piepan.Audio.New`'s argument
    - Remove all plugins; piepan is Lua only
- 0.7.0 (2015-04-08)
    - Add additional Lua support via [gopher-lua](https://github.com/yuin/gopher-lua)
    - Add access token flag
    - Non-script-invoked disconnections are reported though the exit status
    - Remove `-servername` flag (`-lock` + `-insecure` should be used instead)
    - JavaScript plugin: `piepan.Users` and `piepan.Channels` are no longer mapped using string keys ([otto](https://github.com/robertkrimen/otto) needs to be updated before building)
- 0.6.0 (2015-02-11)
    - Fixes due to gumble API changes (see the [gumble API](https://godoc.org/layeh.com/gumble/gumble) if your scripts are not working).
    - Fix crash if `piepan.Process.New` executable did not exist
- 0.5.0 (2015-02-08)
    - Moved to plugin-based system
    - Add certificate locking
    - Add bitrate, volume functions to piepan.Audio
    - Add auto bitrate setting
    - Add Lua plugin
    - Add `piepan.File` to JavaScript plugin
    - Voice targeting is more like the gumble API
- 0.4.0 (2014-12-11)
    - Moved to Go (+ gumble)
    - API has been overhauled. There is no backwards capability with previous versions of piepan
    - JavaScript is now being used as the scripting language
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

## License

MPL 2.0

---

Author: Tim Cooper <<tim.cooper@layeh.com>>
