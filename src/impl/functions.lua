--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.stopAudio()
    if not currentAudio then
        return
    end

    native.stopAudio(currentAudio.ptr)
end

function piepan.disconnect()
    native.disconnect()
end
