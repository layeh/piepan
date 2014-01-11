--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--
--

-- TODO:  coroutines for fetching/returning data we do not yet have (comment/
--        texture hashes)
-- TODO:  kill any timers or threads owned by a script when it reloads

-- Global data
piepan.User.__index = piepan.User
piepan.Message.__index = piepan.Message
piepan.Channel.__index = piepan.Channel
piepan.Timer.__index = piepan.Timer

piepan.server = {}
piepan.args = {}
piepan.scripts = {}
piepan.users = {}
piepan.channels = {}
piepan.Thread.threads = {}
piepan.meta = {}

-- Local data
local functionLock = false
local hasSynced = false
local localUsers = {} -- table of users with the user's session ID as the key
local timers = {}

local native = {
    User = {
        kick = piepan.User.kick,
        moveTo = piepan.User.moveTo,
        ban = piepan.User.ban,
        send = piepan.User.send
    },
    Channel = {
        send = piepan.Channel.send
    },
    Timer = {
        new = piepan.Timer.new,
        cancel = piepan.Timer.cancel
    },
    Thread = {
        new = piepan.Thread.new
    },
    disconnect = piepan.disconnect
}

--
-- Script manager
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
            status, message = pcall(func, arg)
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

--
-- piepan.meta
--
function piepan.meta.__index(table, key)
    return piepan[key]
end

--
-- Timer
--
function piepan.Timer.new(func, timeout, data)
    timeout = math.floor(tonumber(timeout))
    assert(type(func) == "function", "func must be a function")
    assert(timeout > 0 and timeout <= 3600, "timeout is out of range")

    local id = #timers + 1
    local timer = {
        func = func,
        data = data
    }
    timers[id] = timer
    native.Timer.new(id, timeout, timer)
    if timer.ev_timer == nil then
        return nil
    end

    local timerObj = {
        id = id
    }
    setmetatable(timerObj, piepan.Timer)
    return timerObj
end

function piepan.Timer:cancel()
    assert(self ~= nil, "self cannot be nil")

    local timer = timers[self.id]
    if timer == nil then
        return
    end
    native.Timer.cancel(timer.ev_timer)
    timers[self.id] = nil
    self.id = nil
end

function piepan._implOnUserTimer(id)
    assert(functionLock == false, "cannot call implementation functions")

    local timer = timers[id]
    if timer == nil then
        return
    end
    timers[id] = nil
    native.Timer.cancel(timer.ev_timer)

    functionLock = true
    status, message = pcall(timer.func, timer.data)
    if not status then
        print ("Error: timer tick: " .. message)
    end
    functionLock = false
end

--
-- Thread
--
function piepan.Thread.new(worker, callback, data)
    assert(type(worker) == "function", "worker needs to be a function")
    assert(callback == nil or type(callback) == "function",
        "callback needs to be a function or nil")

    local id = #piepan.Thread.threads + 1
    local thread = {
        worker = worker,
        callback = callback,
        data = data
    }
    piepan.Thread.threads[id] = thread
    native.Thread.new(thread, id)
end

function piepan.Thread._implExecute(id)
    local thread = piepan.Thread.threads[id]
    if thread == nil then
        return
    end
    status, val = pcall(thread.worker, thread.data)
    if status == true then
        thread.rtn = val
    end
end

function piepan.Thread._implFinish(id)
    local thread = piepan.Thread.threads[id]
    if thread == nil then
        return
    end
    if thread.callback ~= nil and type(thread.callback) == "function" then
        status, message = pcall(thread.callback, thread.rtn)
        if not status then
            print ("Error: piepan.Thread.finish: " .. message)
        end
    end
    piepan.Thread.threads[id] = nil
end

--
-- User
--

function piepan.User:moveTo(channel)
    assert(self ~= nil, "self cannot be nil")
    assert(getmetatable(channel) == piepan.Channel,
            "channel must be a piepan.Channel")

    if channel == self.channel then
        return
    end
    native.User.moveTo(self, channel.id)
end

function piepan.User:kick(message)
    assert(self ~= nil, "self cannot be nil")

    native.User.kick(self, tostring(message))
end

function piepan.User:ban(message)
    assert(self ~= nil, "self cannot be nil")

    native.User.ban(self, tostring(message))
end

function piepan.User:send(message)
    assert(self ~= nil, "self cannot be nil")

    native.User.send(self, tostring(message))
end

--
-- Channel
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

function piepan.Channel:send(message)
    assert(self ~= nil, "self cannot be nil")

    native.Channel.send(self, tostring(message))
end

--
-- Events
--

function piepan._implOnServerConfig(obj)
    assert(functionLock == false, "cannot call implementation functions")

    if obj.allowHtml ~= nil then
        piepan.server.allowHtml = obj.allowHtml
    end

    piepan._implCall("onConnect")
end

function piepan._implOnServerSync(obj)
    assert(functionLock == false, "cannot call implementation functions")

    piepan.me = localUsers[obj.session]
    if obj.welcomeText ~= nil then
        piepan.server.welcomeText = obj.welcomeText
    end
    hasSynced = true
end

function piepan._implOnMessage(obj)
    assert(functionLock == false, "cannot call implementation functions")

    local message = {
        text = obj.message
    }
    setmetatable(message, piepan.Message)
    if obj.actor ~= nil then
        message.user = localUsers[obj.actor]
    end
    if obj.channels ~= nil then
        -- TODO:  add __len
        message.channels = {}
        for _,v in pairs(obj.channels) do
            message.channels[v] = piepan.channels[v]
        end
    end
    if obj.users ~= nil then
        -- TODO:  add __len
        message.users = {}
        for _,v in pairs(obj.users) do
            local user = localUsers[v]
            if user ~= nil then
                message.users[user.name] = user
            end
        end
    end

    piepan._implCall("onMessage", message)
end

function piepan._implOnUserChange(obj)
    assert(functionLock == false, "cannot call implementation functions")

    local user
    local event = {}
    if localUsers[obj.session] == nil then
        if obj.name == nil then
            return
        end
        user = {
            session = obj.session,
            channel = piepan.channels[0]
        }
        localUsers[obj.session] = user
        piepan.users[obj.name] = user
        setmetatable(user, piepan.User)
        event.isConnected = true
    else
        user = localUsers[obj.session]
    end
    event.user = user

    --
    -- TODO: clear hash if data comes in, and vice-versa.  also needs to resume
    --       the coroutine if we were waiting for that data
    --
    if obj.userId ~= nil then
        user.userId = obj.userId
    end
    if obj.name ~= nil then
        user.name = obj.name
    end
    if obj.channelId ~= nil then
        user.channel = piepan.channels[obj.channelId]
        event.isChangedChannel = true
    end
    if obj.comment ~= nil then
        user.comment = obj.comment
        event.isChangedComment = true
    end
    if obj.isServerMuted ~= nil then
        user.isServerMuted = obj.isServerMuted
    end
    if obj.isServerDeafened ~= nil then
        user.isServerDeafened = obj.isServerDeafened
    end
    if obj.isSelfMuted ~= nil then
        user.isSelfMuted = obj.isSelfMuted
    end
    if obj.isRecording ~= nil then
        user.isRecording = obj.isRecording
    end
    if obj.isSelfDeafened ~= nil then
        user.isSelfDeafened = obj.isSelfDeafened
    end

    if hasSynced then
        piepan._implCall("onUserChange", event)
    end
end

function piepan._implOnUserRemove(obj)
    assert(functionLock == false, "cannot call implementation functions")

    local event = {}
    if localUsers[obj.session] ~= nil then
        -- TODO:  remove reference from Channel -> User?
        local name = localUsers[obj.session].name
        if name ~= nil and piepan.users[name] ~= nil then
            piepan.users[name] = nil
        end
        event.user = localUsers[obj.session]
        localUsers[obj.session] = nil
    end

    if hasSynced and event.user ~= nil then
        event.isDisconnected = true
        piepan._implCall("onUserChange", event)
    end
end

function piepan._implOnChannelRemove(obj)
    assert(functionLock == false, "cannot call implementation functions")

    local channel = piepan.channels[obj.channelId]
    local event = {}
    if channel == nil then
        return
    end
    event.channel = channel

    if channel.parent ~= nil then
        channel.parent.children[channel.id] = nil
        if channel.name ~= nil then
            channel.parent.children[channel.name] = nil
        end
    end
    for k in pairs(channel.children) do
        if k ~= nil then
            k.parent = nil
        end
    end
    piepan.channels[channel.id] = nil

    if hasSynced then
        event.isRemoved = true
        piepan._implCall("onChannelChange", event)
    end
end

function piepan._implOnChannelState(obj)
    assert(functionLock == false, "cannot call implementation functions")

    local channel
    local event = {}
    if piepan.channels[obj.channelId] == nil then
        channel = {
            id = obj.channelId,
            children = {},
            temporary = false,
            users = {}
        }
        piepan.channels[obj.channelId] = channel
        setmetatable(channel, piepan.Channel)
        event.isCreated = true
    else
        channel = piepan.channels[obj.channelId]
    end
    event.channel = channel

    if obj.temporary ~= nil then
        channel.isTemporary = obj.temporary
    end
    if obj.description ~= nil then
        channel.description = obj.description
        event.isChangedDescription = true
    end
    if obj.parentId ~= nil then
        -- Channel got a new parent
        if channel.parent ~= nil and channel.parent.id ~= obj.parentId then
            channel.parent.children[channel.id] = nil
            if channel.name ~= nil then
                channel.parent.children[channel.name] = nil
            end
        end

        channel.parent = piepan.channels[obj.parentId]

        if channel.parent ~= nil then
            channel.parent.children[channel.id] = channel
            if channel.name ~= nil then
                channel.parent.children[channel.name] = channel
            end
        end
        event.isMoved = true
    end
    if obj.name ~= nil then
        if channel.parent ~= nil then
            if channel.name ~= nil then
                channel.parent.children[channel.name] = channel
            end
            channel.parent.children[obj.name] = channel
        end
        channel.name = obj.name
        event.isChangedName = true
    end

    if hasSynced then
        piepan._implCall("onChannelChange", event)
    end
end

function piepan._implOnDisconnect(obj)
    piepan._implCall("onDisconnect", event)
end

--
-- Functions
--
function piepan.disconnect()
    native.disconnect()
end
