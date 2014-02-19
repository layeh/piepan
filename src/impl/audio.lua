--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.Audio.stop()
    if not piepan.Audio.isPlaying() then
        return
    end

    piepan.internal.api.audioStop(piepan.internal.currentAudio.ptr)
end

function piepan.Audio.isPlaying()
    return piepan.internal.currentAudio ~= nil
end
