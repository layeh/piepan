--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

-- TODO:  coroutines for fetching/returning data we do not yet have (comment/
--        texture hashes)
-- TODO:  kill any timers, threads, callbacks owned by a script when it reloads

piepan = {
    User = {},
    UserChange = {},
    Message = {},
    Channel = {},
    ChannelChange = {},
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
        -- table of Users with the user's session ID as the key
        users = {},
        currentAudio
    },
    -- arguments passed to the piepan executable
    args = {},
    scripts = {},
    users = {},
    channels = {},
    threads = {}, -- TODO:  move to internal?
    meta = {},
    timers = {} -- TODO:  move to internal?
}

piepan.User.__index = piepan.User
piepan.UserChange.__index = piepan.UserChange
piepan.Message.__index = piepan.Message
piepan.Channel.__index = piepan.Channel
piepan.ChannelChange.__index = piepan.ChannelChange
piepan.Timer.__index = piepan.Timer
