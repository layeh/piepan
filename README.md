# piepan: an easy to use framework for writing scriptable [Mumble](http://mumble.sourceforge.net/) bots

## Usage

    piepan v0.7.0
    usage: piepan [options] [script files]
    an easy to use framework for writing scriptable Mumble bots
      -access-token=[]: server access token (can be defined multiple times)
      -certificate="": user certificate file (PEM)
      -ffmpeg="ffmpeg": ffmpeg-capable executable for media streaming
      -insecure=false: skip certificate checking
      -key="": user certificate key file (PEM)
      -lock="": server certificate lock file
      -password="": user password
      -server="localhost:64738": address of the server
      -username="piepan-bot": username of the bot

    Script files are defined in the following way:
      [type:[environment:]]filename
        filename: path to script file
        type: type of script file (default: file extension)
        environment: name of environment where script will be executed (default: type)

    Enabled script types:
      Type         Name
      go.lua       Lua (Go)
      lua          Lua (C)
      js           JavaScript

## Scripting documentation

- [JavaScript](https://github.com/layeh/piepan/blob/master/plugins/javascript/README.md)
- [Go Lua](https://github.com/layeh/piepan/blob/master/plugins/golua/README.md)
- [C Lua](https://github.com/layeh/piepan/blob/master/plugins/lua/README.md)

## Building

### Windows

1. Install dependencies
    1. Base dependencies
        1. Download and run [Go installer](https://golang.org/dl/)
        2. Download and run [MSYS2 installer](http://sourceforge.net/projects/msys2/)
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
    - Base package
        - `go get -tags nopkgconfig -u github.com/layeh/piepan`
    - JavaScript plugin (Optional)
        - `go get -u github.com/layeh/piepan/plugins/javascript`
    - Go Lua plugin (Optional)
        - `go get -u github.com/layeh/piepan/plugins/golua`
    - C Lua plugin (Optional)
        - Unavailable on Windows
5. Build piepan
    - `go build -o piepan.exe $GOPATH/src/github.com/layeh/piepan/cmd/piepan/{javascript,golua,main}.go`
6. Run piepan
    - `./piepan.exe ...`

### Ubuntu 14.04

1. Install dependencies
  1. Base dependencies
      - `sudo apt-get install -y git libopus-dev wget python-software-properties software-properties-common pkg-config gcc libav-tools`
  2. Go
      - `sudo add-apt-repository -y ppa:evarlast/golang1.4`
      - `sudo apt-get update`
      - `sudo apt-get install -y golang`
  3. Lua 5.1 (Optional, used for C Lua)
      - `sudo apt-get install -y liblua5.1-0-dev`
2. Create a GOPATH (skip if you already have a GOPATH you want to use)
    - `export GOPATH=$(mktemp -d)`
3. Fetch piepan
    - Base package
        - `go get -u github.com/layeh/piepan`
    - JavaScript plugin (Optional)
        - `go get -u github.com/layeh/piepan/plugins/javascript`
    - Go Lua plugin (Optional)
        - `go get -u github.com/layeh/piepan/plugins/golua`
    - C Lua plugin (Optional)
        - `go get -u github.com/layeh/piepan/plugins/lua`
4. Build piepan (plugins can be removed if they are not wanted)
    - `go build -o piepan $GOPATH/src/github.com/layeh/piepan/cmd/piepan/{javascript,golua,lua,main}.go`
5. Run piepan using `avconv`
    - `./piepan -ffmpeg=avconv ...`

## Changelog

- Next
    - Add additional Lua support via [gopher-lua](https://github.com/yuin/gopher-lua)
    - Add access token flag
    - Non-script-invoked disconnections are reported though the exit status
    - Remove `-servername` flag (`-lock` + `-insecure` should be used instead)
    - JavaScript plugin: `piepan.Users` and `piepan.Channels` are no longer mapped using string keys ([otto](https://github.com/robertkrimen/otto) needs to be updated before building)
- 0.6.0 (2015-02-11)
    - Fixes due to gumble API changes (see the [gumble API](https://godoc.org/github.com/layeh/gumble/gumble) if your scripts are not working).
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

This software is released under the MIT license (see LICENSE).

---

Author: Tim Cooper <<tim.cooper@layeh.com>>
