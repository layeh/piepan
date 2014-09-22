--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

setmetatable(piepan.channels, {
    __call = function (self, path)
        if piepan.channels[0] == nil then
            return nil
        end
        return piepan.channels[0](path)
    end
})

function piepan.Channel:__call(path)
    assert(self ~= nil, "self cannot be nil")

    if path == nil then
        return self
    end
    local channel = self
    for k in path:gmatch("([^/]+)") do
        local current
        if k == "." then
            current = channel
        elseif k == ".." then
            current = channel.parent
        else
            current = channel.children[k]
        end

        if current == nil then
            return nil
        end
        channel = current
    end
    return channel
end

function piepan.Channel:play(info, callback, data)
    assert(self ~= nil, "self cannot be nil")
    assert(type(info) == "string" or type(info) == "table", "info must be a string or table")

    if piepan.internal.currentAudio ~= nil then
        return false
    end

    local filename
    local volume = 1.0
    if type(info) == "table" then
        filename = info.filename
        if info.volume ~= nil then
            volume = info.volume
        end
    else
        filename = info
    end

    local ptr = piepan.internal.api.channelPlay(piepan.internal.state,
        piepan.internal.opus.encoder, filename, volume)
    if not ptr then
        return false
    end
    piepan.internal.currentAudio = {
        callback = callback,
        callbackData = data,
        ptr = ptr
    }
    return true
end

function piepan.internal.events.onAudioFinished()
    assert (piepan.internal.currentAudio ~= nil, "audio must be playing")

    local callback = piepan.internal.currentAudio.callback

    if type(callback) == "function" then
        local data = piepan.internal.currentAudio.callbackData
        piepan.internal.currentAudio = nil
        piepan.internal.runCallback(callback, data)
    else
        piepan.internal.currentAudio = nil
    end
end

function piepan.Channel:send(message)
    assert(self ~= nil, "self cannot be nil")

    piepan.internal.api.channelSend(self, tostring(message))
end

function piepan.Channel:setDescription(description)
    assert(self ~= nil, "self cannot be nil")
    assert(type(description) == "string" or description == nil,
        "description must be a string or nil")

    if description == nil then
        description = ""
    end
    piepan.internal.api.channelSetDescription(self, description)
end

function piepan.Channel:remove()
    assert(self ~= nil, "self cannot be nil")

    piepan.internal.api.channelRemove(self)
end

function piepan.Channel:resolveHashes()
    assert(self ~= nil, "self cannot be nil")

    local request
    if self.descriptionHash == nil then
        return
    end

    local running = coroutine.running()
    if piepan.internal.resolving.channels[self.id] == nil then
        piepan.internal.resolving.channels[self.id] = {running}
        request = true
    else
        if #piepan.internal.resolving.channels[self.id] <= 0 then
            request = true
        end
        table.insert(piepan.internal.resolving.channels[self.id], running)
    end
    if request then
        piepan.internal.api.resolveHashes(nil, nil, {self.id})
    end
    coroutine.yield()
end
