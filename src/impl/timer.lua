--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.Timer.new(func, timeout)
    assert(type(func) == "function", "func must be a function")
    assert(type(timeout) == "number" and timeout > 0 and timeout <= 3600,
        "timeout is out of range")

    local id = #piepan.internal.timers + 1
    local timerObj = {
        id = id
    }
    piepan.internal.timers[id] = {
        func = func,
        ptr = nil,
        state = nil
    }
    piepan.internal.api.timerNew(piepan.internal.timers[id], id, timeout,
        piepan.internal.state)

    setmetatable(timerObj, piepan.Timer)
    return timerObj
end

function piepan.Timer:cancel()
    assert(self ~= nil, "self cannot be nil")

    local timer = piepan.internal.timers[self.id]
    if timer == nil then
        return
    end
    piepan.internal.api.timerCancel(timer.ptr)
    piepan.internal.timers[self.id] = nil
    self.id = nil
end

function piepan.internal.events.onUserTimer(id)
    local timer = piepan.internal.timers[id]
    if timer == nil then
        return
    end

    piepan.internal.timers[id] = nil
    piepan.internal.runCallback(timer.func)
end
