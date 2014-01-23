--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.User:moveTo(channel)
    assert(self ~= nil, "self cannot be nil")
    assert(getmetatable(channel) == piepan.Channel,
            "channel must be a piepan.Channel")

    if channel == self.channel then
        return
    end
    native.User.moveTo(self, channel.id)
end

function piepan.User:kick(message)
    assert(self ~= nil, "self cannot be nil")

    native.User.kick(self, tostring(message))
end

function piepan.User:ban(message)
    assert(self ~= nil, "self cannot be nil")

    native.User.ban(self, tostring(message))
end

function piepan.User:send(message)
    assert(self ~= nil, "self cannot be nil")

    native.User.send(self, tostring(message))
end
