--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.internal.events.onLoadScript(argument, ptr)
    local index
    local entry

    if type(argument) == "string" then
        index = #piepan.scripts + 1
        entry = {
            filename = argument,
            ptr = ptr,
            environment = {
                print = print,
                assert = assert,
                collectgarbage = collectgarbage,
                dofile = dofile,
                error = error,
                getmetatable = getmetatable,
                ipairs = ipairs,
                load = load,
                loadfile = loadfile,
                next = next,
                pairs = pairs,
                pcall = pcall,
                print = print,
                rawequal = rawequal,
                rawget = rawget,
                rawlen = rawlen,
                rawset = rawset,
                require = require,
                select = select,
                setmetatable = setmetatable,
                tonumber = tonumber,
                tostring = tostring,
                type = type,
                xpcall = xpcall,

                bit32 = bit32,
                coroutine = coroutine,
                debug = debug,
                io = io,
                math = math,
                os = os,
                package = package,
                string = string,
                table = table
            }
        }
    elseif type(argument) == "number" then
        index = argument
        entry = piepan.scripts[index]
    else
        return false, "invalid argument"
    end

    local script, message = loadfile(entry.filename, "bt", entry.environment)
    if script == nil then
        return false, message
    end
    entry.environment.piepan = {}
    local status, message = pcall(script)
    if status == false then
        return false, message
    end

    piepan.scripts[index] = entry
    if type(entry.environment.piepan) == "table" then
        setmetatable(entry.environment.piepan, piepan.internal.meta)
    end

    return true, index, ptr
end

--
-- Callback execution
--
function piepan.internal.triggerEvent(name, ...)
    for _,script in pairs(piepan.scripts) do
        local func = rawget(script.environment.piepan, name)
        if type(func) == "function" then
            piepan.internal.runCallback(func, ...)
        end
    end
end

function piepan.internal.runCallback(func, ...)
    assert(type(func) == "thread" or type(func) == "function",
        "func should be a coroutine or a function")

    local routine
    if type(func) == "thread" then
        routine = func
    else
        routine = coroutine.create(func)
    end
    local status, message = coroutine.resume(routine, ...)
    if not status then
        print ("Error: " .. message)
    end
end

--
-- Argument parsing
--
function piepan.internal.events.onArgument(key, value)
    assert(type(key) ~= nil, "key cannot be nil")

    value = value or ""
    if piepan.args[key] == nil then
        piepan.args[key] = {value}
    else
        table.insert(piepan.args[key], value)
    end
end
