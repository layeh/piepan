--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
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
    if obj.maxBandwidth ~= nil then
        piepan.server.maxBandwidth = obj.maxBandwidth
    end
    piepan.server.synced = true
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

    if piepan.server.synced then
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

    if piepan.server.synced and event.user ~= nil then
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

    if piepan.server.synced then
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

    if piepan.server.synced then
        piepan._implCall("onChannelChange", event)
    end
end

function piepan._implOnDisconnect(obj)
    piepan._implCall("onDisconnect", event)
end
