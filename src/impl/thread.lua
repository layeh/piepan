--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.Thread.new(worker, callback, data)
    assert(type(worker) == "function", "worker needs to be a function")
    assert(callback == nil or type(callback) == "function",
        "callback needs to be a function or nil")

    local id = #piepan.internal.threads + 1
    local thread = {
        worker = worker,
        callback = callback,
        data = data
    }
    piepan.internal.threads[id] = thread
    piepan.internal.api.threadNew(thread, id)
end

function piepan.internal.events.onThreadExecute(id)
    local thread = piepan.internal.threads[id]
    if thread == nil then
        return
    end
    status, val = pcall(thread.worker, thread.data)
    if status == true then
        thread.rtn = val
    end
end

function piepan.internal.events.onThreadFinish(id)
    local thread = piepan.internal.threads[id]
    if thread == nil then
        return
    end
    if thread.callback ~= nil and type(thread.callback) == "function" then
        status, message = pcall(thread.callback, thread.rtn)
        if not status then
            print ("Error: piepan.Thread.finish: " .. message)
        end
    end
    piepan.internal.threads[id] = nil
end
