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

function piepan.User:register()
    assert(self ~= nil, "self cannot be nil")

    piepan.internal.api.userRegister(self)
end

function piepan.User:resolveHashes()
    assert(self ~= nil, "self cannot be nil")
    local comment, texture
    local request
    local count = 0

    if self.textureHash ~= nil then
        texture = {self.session}
        count = count + 1
    end
    if self.commentHash ~= nil then
        comment = {self.session}
        count = count + 1
    end
    if texture == nil and comment == nil then
        return
    end

    local running = coroutine.running()
    local tbl = {
        routine = running,
        count = count
    }
    if piepan.internal.resolving.users[self.session] == nil then
        piepan.internal.resolving.users[self.session] = {tbl}
        request = true
    else
        if #piepan.internal.resolving.users <= 0 then
            request = true
        end
        table.insert(piepan.internal.resolving.users[self.session], tbl)
    end
    if request then
        piepan.internal.api.resolveHashes(texture, comment, nil)
    end
    coroutine.yield()
end

function piepan.User:setTexture(bytes)
    assert(self ~= nil, "self cannot be nil")
    assert(type(bytes) == "string" or bytes == nil, "bytes must be a string or nil")

    if bytes == nil then
        bytes = ""
    end

    piepan.internal.api.userSetTexture(bytes)
end
