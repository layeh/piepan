--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.Timer.new(func, timeout, data)
    assert(type(func) == "function", "func must be a function")
    assert(type(timeout) == "number" and timeout > 0 and timeout <= 3600,
        "timeout is out of range")

    local id = #piepan.timers + 1
    local timerObj = {
        id = id
    }
    local timer = {
        func = func,
        data = data,
        handle = timerObj,
        ptr = piepan.internal.api.timerNew(id, timeout)
    }
    piepan.timers[id] = timer

    setmetatable(timerObj, piepan.Timer)
    return timerObj
end

function piepan.Timer:cancel()
    assert(self ~= nil, "self cannot be nil")

    local timer = piepan.timers[self.id]
    if timer == nil then
        return
    end
    piepan.internal.api.timerCancel(timer.ptr)
    piepan.timers[self.id] = nil
    self.id = nil
end

function piepan.internal.events.onUserTimer(id)
    local timer = piepan.timers[id]
    if timer == nil then
        return
    end
    piepan.timers[id] = nil

    status, message = pcall(timer.func, timer.data)
    if not status then
        print ("Error: timer tick: " .. message)
    end
end
