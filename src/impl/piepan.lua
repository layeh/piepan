--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

-- TODO:  kill any timers, threads, callbacks owned by a script when it reloads

piepan = {
    Audio = {},
    User = {},
    UserChange = {},
    Message = {},
    Channel = {},
    ChannelChange = {},
    PermissionDenied = {},
    Permissions = {},
    Thread = {},
    Timer = {},

    server = {
        -- has the client been fully synced with the server yet?
        synced = false
    },
    internal = {
        api = {},
        opus = {},
        events = {},
        threads = {},
        timers = {},
        meta = {},
        -- table of Users with the user's session ID as the key
        users = {},
        permissionsMap = {
            write   = 0x1,
            traverse = 0x2,
            enter = 0x4,
            speak = 0x8,
            muteDeafen = 0x10,
            move = 0x20,
            makeChannel = 0x40,
            linkChannel = 0x80,
            whisper = 0x100,
            textMessage = 0x200,
            makeTemporaryChannel = 0x400,
            kick = 0x10000,
            ban = 0x20000,
            register = 0x40000,
            registerSelf = 0x80000
        },
        resolving = {
            users = {},
            channels = {}
        },
        currentAudio,
        state
    },
    -- arguments passed to the piepan executable
    args = {},
    scripts = {},
    users = {},
    channels = {}
}

piepan.Audio.__index = piepan.Audio
piepan.User.__index = piepan.User
piepan.UserChange.__index = piepan.UserChange
piepan.Message.__index = piepan.Message
piepan.Channel.__index = piepan.Channel
piepan.ChannelChange.__index = piepan.ChannelChange
piepan.PermissionDenied.__index = piepan.PermissionDenied
piepan.Permissions.__index = piepan.Permissions
piepan.Timer.__index = piepan.Timer
