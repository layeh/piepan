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

    if currentAudio ~= nil then
        return false
    end

    local ptr = native.Channel.play(piepan.internal.opus.encoder, filename)
    if not ptr then
        return false
    end
    currentAudio = {
        callback = callback,
        callbackData = data,
        ptr = ptr
    }
    return true
end

function piepan.Channel._implAudioFinished()
    assert (currentAudio ~= nil, "audio must be playing")

    if type(currentAudio.callback) == "function" then
        status, message = pcall(currentAudio.callback, currentAudio.callbackData)
        if not status then
            print ("Error: " .. message)
        end
    end

    currentAudio = nil
end

function piepan.Channel:send(message)
    assert(self ~= nil, "self cannot be nil")

    native.Channel.send(self, tostring(message))
end
