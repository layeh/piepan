# piepan: an easy to use framework for writing scriptable [Mumble](http://mumble.sourceforge.net/) bots

## Usage

    usage: piepan [options] [script files]
    an easy to use framework for writing scriptable Mumble bots
      -certificate="": user certificate file (PEM)
      -insecure=false: skip certificate checking
      -key="": user certificate key file (PEM)
      -lock="": server certificate lock file
      -password="": user password
      -server="localhost:64738": address of the server
      -servername="": override server name used in TLS handshake
      -username="piepan-bot": username of the bot

    Script files are defined in the following way:
      [type:[environment:]]filename
        filename: path to script file
        type: type of script file (default: file extension)
        environment: name of environment where script will be executed (default: type)

    Enabled script types:
      js

## Scripting documentation

- [JavaScript](https://github.com/layeh/piepan/blob/master/plugins/javascript/README.md)

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

## Changelog

- Next
    - Moved to plugin-based system
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

## License

This software is released under the MIT license (see LICENSE).

---

Author: Tim Cooper <<tim.cooper@layeh.com>>
