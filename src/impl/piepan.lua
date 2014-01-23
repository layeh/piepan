--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

-- TODO:  coroutines for fetching/returning data we do not yet have (comment/
--        texture hashes)
-- TODO:  kill any timers, threads, callbacks owned by a script when it reloads

-- Global data
piepan.User.__index = piepan.User
piepan.Message.__index = piepan.Message
piepan.Channel.__index = piepan.Channel
piepan.Timer.__index = piepan.Timer

piepan.server = {
    synced = false
}
piepan.internal = {
    opus = {
        encoder = nil
    }
}
piepan.args = {}
piepan.scripts = {}
piepan.users = {}
piepan.channels = {}
piepan.threads = {}
piepan.meta = {}
piepan.timers = {}

-- Local data
local functionLock = false
local localUsers = {} -- table of users with the user's session ID as the key
local currentAudio

local native = {
    User = {
        kick = piepan.User.kick,
        moveTo = piepan.User.moveTo,
        ban = piepan.User.ban,
        send = piepan.User.send
    },
    Channel = {
        play = piepan.Channel.play,
        send = piepan.Channel.send
    },
    Timer = {
        new = piepan.Timer.new,
        cancel = piepan.Timer.cancel
    },
    Thread = {
        new = piepan.Thread.new
    },
    stopAudio = piepan.stopAudio,
    disconnect = piepan.disconnect
}
