--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan._implLoadScript(argument)
    assert(functionLock == false, "cannot call implementation functions")

    local index
    local entry

    if type(argument) == "string" then
        index = #piepan.scripts + 1
        entry = {
            filename = argument,
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
        setmetatable(entry.environment.piepan, piepan.meta)
    end

    return true, index
end

function piepan._implCall(name, arg)
    assert(type(name) == "string", "name must be a string")

    functionLock = true
    for _,script in pairs(piepan.scripts) do
        local func = rawget(script.environment.piepan, name)
        if type(func) == "function" then
            local status, message = pcall(func, arg)
            if not status then
                print ("Error: " .. message)
            end
        end
    end
    functionLock = false
end

--
-- Argument parsing
--
function piepan._implArgument(key, value)
    assert(type(key) ~= nil, "key cannot be nil")

    value = value or ""
    if piepan.args[key] == nil then
        piepan.args[key] = {value}
    else
        table.insert(piepan.args[key], value)
    end
end
