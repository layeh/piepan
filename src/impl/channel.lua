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

function piepan.Channel:play(filename, callback, data)
    assert(self ~= nil, "self cannot be nil")
    assert(type(filename) == "string", "filename must be a string")

    if piepan.internal.currentAudio ~= nil then
        return false
    end

    local ptr = piepan.internal.api.channelPlay(piepan.internal.opus.encoder,
        filename)
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

    if type(piepan.internal.currentAudio.callback) == "function" then
        status, message = pcall(piepan.internal.currentAudio.callback,
            piepan.internal.currentAudio.callbackData)
        if not status then
            print ("Error: " .. message)
        end
    end

    piepan.internal.currentAudio = nil
end

function piepan.Channel:send(message)
    assert(self ~= nil, "self cannot be nil")

    piepan.internal.api.channelSend(self, tostring(message))
end
