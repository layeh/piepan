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
    piepan.internal.api.userMoveTo(self, channel.id)
end

function piepan.User:kick(message)
    assert(self ~= nil, "self cannot be nil")

    piepan.internal.api.userKick(self, tostring(message))
end

function piepan.User:ban(message)
    assert(self ~= nil, "self cannot be nil")

    piepan.internal.api.userBan(self, tostring(message))
end

function piepan.User:send(message)
    assert(self ~= nil, "self cannot be nil")

    piepan.internal.api.userSend(self, tostring(message))
end

function piepan.User:setComment(comment)
    assert(self ~= nil, "self cannot be nil")
    assert(type(comment) == "string" or comment == nil,
        "comment must be a string or nil")

    if comment == nil then
        comment = ""
    end
    piepan.internal.api.userSetComment(self, comment)
end
