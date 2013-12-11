function piepan.onConnect()
    print ("echo loaded!")
end

function piepan.onMessage(message)
    if message.user == nil then
        return
    end
    piepan.me.channel:send(message.text)
end
