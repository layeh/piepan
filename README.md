# piepan: a bot framework for Mumble

piepan is an easy to use bot framework for interacting with a [Mumble](http://mumble.sourceforge.net/) server.  Here is a simple script that will echo back any chat message that is sent to it:

    -- echo.lua
    function piepan.onMessage(msg)
        piepan.me.channel:send(msg.text)
    end

The above script can be started from the command line:

    # piepan echo.lua

## Usage

    usage: piepan [options] [scripts...]
    a bot framework for Mumble

      -u <username>       username of the bot (has no effect if the certificate
                          has been registered with the server under a different
                          name)
      -s <server>         address of the server (default: localhost)
      -p <port>           port of the server (default: 64738)
      -pw <file>          read server password from the given file (when file is -,
                          standard input will be read)
      -t <file>           read access tokens (one per line) from the given file
      -c <certificate>    certificate to use for the connection
      -k <keyfile>        key file to use for the connection (defaults to the
                          certificate file)
      -d                  enable development mode, which automatically reloads
                          scripts when they are modified
      --<name>[=<value>]  a key-value pair that will be accessible from the scripts
      -h                  display this help
      -v                  show version


## Programming reference

The following section describes the API that is available for script authors.  Please note that the current API does not contain all of the features that are defined in the Mumble protocol.

### Prototypes

#### `piepan.User`

- **`int session`**: the session ID of the user
- **`int userId`**: the registered user ID of the user
- **`string name`**: the username of the user
- **`string comment`**: the user's comment (FIXME)
- **`bool isServerMuted`**: is the user muted by the server
- **`bool isServerDeafened`**: has the user been deafened by the server
- **`bool isSelfMuted`**: is the user muted by the him/herself
- **`bool isSelfDeafened`**: has the user been deafened by him/herself
- **`bool isRecording`**: is the user recording channel audio
- **`piepan.Channel channel`**: the channel that the user is currently in
- **`void moveTo(self, piepan.Channel channel)`**: moves the user to the given `channel`
- **`void send(self, string message)`**: sends a message to the user
- **`void kick(self [, string reason])`**: kicks the user from the server with an optional reason
- **`void ban(self [, string reason])`**: bans the user from the server with an optional reason

#### `piepan.Message`

- **`string text`**: the message text
- **`piepan.User user`**: the user who sent the message (this can be `nil`)
- **`piepan.Channel channels`**: a table of channels the message was sent to, with the key being the channel ID and the value being the corresponding channel table
- **`piepan.User users`**: a table of users the message was sent to, with the key being the user name and the value being their corresponding user table

#### `piepan.Channel`

- **`int id`**: the unique channel identifier
- **`string name`**: the channel name
- **`string description`**: the description of the channel (FIXME)
- **`piepan.Channel parent`**: the parent channel
- **`bool isTemporary`**: is the channel temporary
- **`void send(self, string message)`**: sends a message to the channel

    Example:

        -- sends a message to the channel the bot is currently in
        piepan.me.channel:send("Hello Everyone!")

- **`piepan.Channel __call(self, string path)`**: returns the child at the end of the path. The path items are separated by slashes (`/`).  The path item `.` refers to the current channel, the item `..` refers to the parent channel, and all other items refer to the child channel.

    Example:

        -- moves user to the sibling channel named test
        local channel = user.channel("../test")
        user:moveTo(channel)

#### `piepan.Timer`

- **`piepan.Timer new(func, timeout [, data])`**: Creates a new timer.  After `timeout` seconds elapses, `func` will be called with `data` as its first and only parameter.

    The (arbitrary) range of `timeout` is [1, 3600] (1 second to 1 hour).

    Once a timer has been fired or canceled, its reference is no longer valid.

- **`void cancel(self)`**: Cancels a timer.

#### `piepan.Thread`

- **`void new(worker [, callback, data])`**: Starts executing the global function `worker` in a new thread, with the argument `data`.

    The worker function should only use local variables.  Any use or modification of global variables is undefined.  Values that this function needs should be passed via the `data` argument.

    An optional `callback` will be executed on the main thread after `worker` completes.  It will be passed the value that `worker` returns.

#### `piepan.UserChange` (hidden)

- **`piepan.User user`**: the user that changed
- **`bool isConnected`**:  if the user connected to the server
- **`bool isDisconnected`**: if the user disconnected from the server
- **`bool isChangedChannel`**:  if the user moved to a new channel
- **`bool isChangedComment`**: if the user's comment changed

#### `piepan.ChannelChange` (hidden)

- **`piepan.Channel channel`**: the channel that was changed
- **`bool isCreated`**: if the channel was created
- **`bool isRemoved`**: if the channel was removed
- **`bool isMoved`**: if the channel moved in the tree
- **`bool isChangedName`**: if the channel name changed
- **`bool isChangedDescription`**: if the channel description changed

### Variables

#### `piepan.users`

Table containing each connected user on the server, with the keys being the name of the user and the value being their corresponding `piepan.User` table.

Example:

    -- prints the usernames of all the users connected to the server to standard
    -- output
    for name,_ in pairs(piepan.users) do
        print (name)
    end

#### `piepan.channels`

Table that contains all of the channels that are on the server. The channels are mapped by their channel IDs. `piepan.channels[0]` is the server's root channel.

`piepan.channels.__call` is mapped to `piepan.channels[0].__call`, therefore the following can be done:

    local channel = piepan.channels("A/B/C")
    piepan.me:moveTo(channel)

#### `piepan.me`

The `piepan.User` table that references yourself.

#### `piepan.server`

Table containing information about the server.  This table may have the fields:

- **`bool allowHtml`**: are HTML messages allowed to be sent to the server
- **`string welcomeText`**: the server's welcome text

#### `piepan.args`

Table that is populated with the command line arguments that are in the form: `--key=value` or `--key` (in the latter case, `value` is an empty string).

The values are stored in the array `piepan.args[key]`, which allows multiple arguments with the same key.

Example:

    -- print a list of all of the admins
    for _,admin in ipairs(piepan.args.admin or {}) do
        print (admin)
    end

    # program execution
    > piepan --admin=user1 --admin=user2 ...

### Functions

#### `piepan.disconnect()`

Disconnects from the server.

### Callbacks

These are functions that can be defined in script files.  They will be called when the corresponding event happens.

#### `piepan.onConnect()`

Called when connection to the server has been made. This is where a script should perform its initialization.

#### `piepan.onDisconnect()`

Called when connection to the server has been lost or after `piepan.disconnect()` is called.

#### `piepan.onMessage(piepan.Message message)`

Called when a text message `message` is received.

#### `piepan.onUserChange(piepan.UserChange event)`

Called when a user's status changes.

#### `piepan.onChannelChange(piepan.ChannelChange event)`

Called when a channel changes state (e.g. is added or removed).

## Changelog

- 0.1 (2013-12-11)
    - Initial Release

## Building

1. Ensure all of the requirements are installed
2. Run `make` inside of the project directory
3. The `piepan` executable should appear in the directory

## Requirements

- [OpenSSL](http://www.openssl.org/)
- [Lua 5.2](http://www.lua.org/)
- [libev](http://libev.schmorp.de/)
- [protobuf-c](https://github.com/protobuf-c/protobuf-c)

## License

This software is released under the MIT license (see LICENSE).

---

Author: Tim Cooper <<tim.cooper@layeh.com>>
