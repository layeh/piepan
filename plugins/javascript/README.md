# JavaScript plugin documentation

## API

The following section describes the API that is available for piepan JavaScript scripts.

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

#### [`Channels`](https://godoc.org/github.com/layeh/gumble/gumble#Channels) `piepan.Channels`

Object that contains all of the channels that are on the server. The channels are mapped by their channel IDs (as a string). `piepan.Channels["0"]` is the server's root channel.

Note: `Channel.Channels` and `Channel.Users` cannot be iterated over.

#### `piepan.Disconnect()`

Disconnects from the server.

#### `piepan.File`

- `piepan.File Open(string filename [, string mode])`: opens `filename` with `mode`. The following modes are supported:
    - `r`: read only
    - `r+`: read-write
    - `w`: write only, create file if it does not exist
    - `w+`: read-write, create file if it does not exist, truncate file
    - `a`: write only, create file if it does not exist, append file
    - `a+`: read-write, create file if it does not exist, append file
- `void Close()`: close the file.
- `string Read([int n])`: reads `n` bytes from the file, starting at the current offset. Omitting or setting `n` less than 1 will read until the end of file.
- `int Seek(int offset [, string whence])`: Seeks to `offset`, relative to `whence`. `whence` can be one of:
    - `set`: `offset` relative to start of file
    - `cur`: `offset` relative to the current offset (default)
    - `end`: `offset` relative to end of file
- `int Write(string data)`: write `data` to file. Returns number of bytes written.

#### `piepan.On(string event, function callback)`

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

#### `piepan.Process`

- `piepan.Process New(function callback, string command, string arguments...)`: Executes `command` in a new process with the given arguments. The function `callback` is executed once the process has completed, passing if the execution was successful and the contents of standard output.

- `void Kill()`: Kills the process.

#### [`User`](https://godoc.org/github.com/layeh/gumble/gumble#User) `piepan.Self`

The `User` object that references yourself.

#### `piepan.Timer`

- `piepan.Timer New(function callback, int timeout)`: Creates a new timer.  After at least `timeout` milliseconds, `callback` will be executed.

- `void Cancel()`: Cancels the timer.

#### [`Users`](https://godoc.org/github.com/layeh/gumble/gumble#Users) `piepan.Users`

Object containing each connected user on the server, with the keys being the session ID of the user (as a string) and the value being their corresponding `piepan.User` object.

Example:

    // Print the names of the connected users to standard output
    for (var k in piepan.Users) {
      var user = piepan.Users[k];
      console.log(user.Name);
    }
