--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.stopAudio()
    if not piepan.internal.currentAudio then
        return
    end

    piepan.internal.api.stopAudio(piepan.internal.currentAudio.ptr)
end

function piepan.disconnect()
    piepan.internal.api.disconnect()
end
