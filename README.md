# [piepan][1]: a bot framework for Mumble

piepan is an easy to use framework for writing scriptable [Mumble](http://mumble.sourceforge.net/) bots using JavaScript.  Here is a simple script that will echo back any chat message that is sent to it:

    -- echo.js
    piepan.On('message', function(e) {
      if (e.Sender == null) {
        return;
      }
      piepan.Self.Channel().Send(e.Message, false);
    });


The above script can be started from the command line:

    $ piepan echo.js

## Usage

    usage: piepan [options] [scripts...]
    a scriptable bot framework for Mumble
      -certificate="": user certificate file (PEM)
      -insecure=false: skip certificate checking
      -key="": user certificate key file (PEM)
      -lock="": server certificate lock file
      -password="": user password
      -server="localhost:64738": address of the server
      -servername="": override server name used in TLS handshake
      -username="piepan-bot": username of the bot

## Building

### Windows

1. Install dependencies (use 32-bit/386/x86 installers)
  1. Download and run [Go installer](https://golang.org/dl/)
  2. Download and run [MinGW installer](http://www.mingw.org/)
    - Configure [PATH variable](http://www.mingw.org/wiki/getting_started#toc7)
    - Install the following packages from the MinGW installation manager:
      - mingw-developer-toolkit
      - mingw32-base
      - msys-base
      - msys-wget
  3. Download [pkg-config-lite](http://sourceforge.net/projects/pkgconfiglite/)
    - Extract zip contents to the MinGW installation folder
  4. Download and run [Git installer](http://git-scm.com/download/win)
    - Select "Use Git from the Windows Command Prompt" during installation
2. Open MSYS terminal (defaults to `C:\MinGW\msys\1.0\msys.bat`)
  1. Download and install Opus
    - `wget http://downloads.xiph.org/releases/opus/opus-1.1.tar.gz`
    - `tar xzf opus-1.1.tar.gz`
    - `(cd opus-1.1; ./configure && make install)`
  2. Configure GOPATH
    - `export GOPATH=$(mktemp -d)`
  3. Configure pkg-config path
    - `export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig/`
  4. Build piepan
    - `go get -u github.com/layeh/piepan/cmd/piepan`
  5. Copy .exe to current directory
    - `cp "$GOPATH/bin/piepan.exe" .`
  6. Run piepan.exe

### Ubuntu 14.04

1. Install dependencies
  2. Base dependencies
    - `sudo apt-get install -y git libopus-dev wget python-software-properties software-properties-common pkg-config gcc`
  2. Golang
    - `sudo add-apt-repository -y ppa:evarlast/golang1.4`
    - `sudo apt-get update`
    - `sudo apt-get install -y golang`
  3. ffmpeg (if you would like to play media files)
    - `sudo add-apt-repository -y ppa:jon-severinsson/ffmpeg`
    - `sudo apt-get update`
    - `sudo apt-get install -y ffmpeg`
2. Create a GOPATH (skip if you already have a GOPATH you want to use)
  - `export GOPATH=$(mktemp -d)`
3. Build piepan
  - `go get -u github.com/layeh/piepan/cmd/piepan`
4. Copy binary to current directory
  - `cp "$GOPATH/bin/piepan" .`

## Programming reference

The following section describes the API that is available for piepan scripts.

piepan is built using the [gumble](https://github.com/layeh/gumble) library. Documentation for types not part of piepan itself (including User and Channel) can be found in the [gumble documentation](https://godoc.org/github.com/layeh/gumble/gumble).

### `piepan.Audio`

- `void Play(object obj)`: Plays the media file `obj.filename`. `obj.callback` can be defined as a function that is called after the playback has completed.
- [`VoiceTarget`](https://godoc.org/github.com/layeh/gumble/gumble#VoiceTarget)   `NewTarget(int id)`: Create a new voice target object.
- `void SetTarget(VoiceTarget target)` sets the target of subsequent `piepan.Audio.Play()` calls. Call this function with no arguments to remove any voice targeting.
- `void Stop()`: Stops the currently playing stream.
- `bool IsPlaying()`: Returns true if an stream is currently playing, false otherwise.
- `int Bitrate()`: Returns the bitrate of the audio encoder.
- `void SetBitrate(int bitrate)`: Sets the bitrate of the audio encoder. Calling this function will override the automatically-configured, optimal bitrate.
- `float Volume()`: Returns the audio volume.
- `void SetVolume(float volume)`: Sets the volume of transmitted audio (default: 1.0).

### [`Channels`](https://godoc.org/github.com/layeh/gumble/gumble#Channels) `piepan.Channels`

Object that contains all of the channels that are on the server. The channels are mapped by their channel IDs (as a string). `piepan.Channels["0"]` is the server's root channel.

Note: [`Channel.Channels()`](https://godoc.org/github.com/layeh/gumble/gumble#Channel.Channels) and [`Channel.Users()`](https://godoc.org/github.com/layeh/gumble/gumble#Channel.Users) cannot be iterated over.

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
- `permissionDenied` (Arguments: [`PermissionDeniedEvent event`](https://godoc.org/github.com/layeh/gumble/gumble#PermissionDeniedEvent))
    - Called when a requested action could not be performed.

Note: events with a `Type` field have slight changes than what is documented in gumble:

1. The `Type` field is changed to a number.
2. Individual bit flag values are added to the event as booleans prefixed with `Is`
    - `DisconnectEvent`
        - `IsError`
        - `IsUser`
        - `IsOther`
        - `IsVersion`
        - `IsUserName`
        - `IsUserCredentials`
        - `IsServerPassword`
        - `IsUsernameInUse`
        - `IsServerFull`
        - `IsNoCertificate`
        - `IsAuthenticatorFail`
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

### [`User`](https://godoc.org/github.com/layeh/gumble/gumble#User) `piepan.Self`

The `User` object that references yourself.

### `piepan.Timer`

- `piepan.Timer New(function callback, int timeout)`: Creates a new timer.  After at least `timeout` milliseconds, `callback` will be executed.

- `void Cancel()`: Cancels the timer.

### [`Users`](https://godoc.org/github.com/layeh/gumble/gumble#Users) `piepan.Users`

Object containing each connected user on the server, with the keys being the session ID of the user (as a string) and the value being their corresponding `piepan.User` object.

Example:

    // Print the names of the connected users to standard output
    for (var k in piepan.Users) {
      var user = piepan.Users[k];
      console.log(user.Name());
    }

## Changelog

- Next
    - Voice targeting is more like the gumble API
    - Add certificate locking
    - Add bitrate, volume functions to piepan.Audio
    - Add auto bitrate setting
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

## Requirements

- [gumble](https://github.com/bontibon/gumble/tree/master/gumble)
- [gumble_ffmpeg](https://github.com/bontibon/gumble/tree/master/gumble_ffmpeg)
- [otto](https://github.com/robertkrimen/otto)

## License

This software is released under the MIT license (see LICENSE).

---

Author: Tim Cooper <<tim.cooper@layeh.com>>

[1]: https://github.com/layeh/piepan
