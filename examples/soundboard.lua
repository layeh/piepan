--
-- Soundboard example.
--
-- Sounds are triggered when #<keyword> appears in a text message, where
-- keyword has been defined below in the sounds table.
--

-- Boolean if users need to be registered on the server to trigger sounds
local require_registered = true

-- Boolean if sounds should stop playing when another is triggered
local interrupt_sounds = false

-- Boolean if the bot should move into the user's channel to play the sound
local should_move = false

-- Table with keys being the keywords and values being the sound files
local sounds = {
    cheer = "cheer.ogg",
    hello = "hello.ogg",
    huh = "huh.ogg",
    image = "image.ogg",
    lol = "lol.ogg",
    mock = "mock.ogg",
    nice = "nice.ogg"
}

-- Sound file path prefix
local prefix = "examples/sounds/"

---------------

function piepan.onConnect()
    if piepan.args.soundboard then
        prefix = piepan.args.soundboard
    end
    print ("Soundboard loaded!")
end

function piepan.onMessage(msg)
    if msg.user == nil then
        return
    end

    local search = string.match(msg.text, "#(%w+)")
    if not search or not sounds[search] then
        return
    end
    local soundFile = prefix .. sounds[search]
    if require_registered and msg.user.userId == nil then
        msg.user:send("You must be registered on the server to trigger sounds.")
        return
    end
    if piepan.Audio.isPlaying() and not interrupt_sounds then
        return
    end
    if piepan.me.channel ~= msg.user.channel then
        if not should_move then
            return
        end
        piepan.me:moveTo(msg.user.channel)
    end

    piepan.Audio.stop()
    piepan.me.channel:play(soundFile)
end
