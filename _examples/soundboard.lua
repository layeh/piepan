--
-- Soundboard example.
--
-- Sounds are triggered when #<keyword> appears in a text message, where
-- keyword has been defined below in the sounds table.
--
do

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
local prefix = os.getenv("PREFIX") or "_examples/sounds/"

---------------

piepan.On("connect", function()
  print("Soundboard loaded!")
end)

piepan.On("message", function(e)
  if e.Sender == nil then
    return
  end

  local search = string.match(e.Message, "#(%w+)")
  if not search or not sounds[search] then
    return
  end
  local soundFile = prefix .. sounds[search]
  if require_registered and e.Sender.UserID == 0 then
    e.Sender:Send("You must be registered on the server to trigger sounds.")
    return
  end
  if piepan.Audio.IsPlaying() and not interrupt_sounds then
    return
  end
  if piepan.Self.Channel ~= e.Sender.Channel then
    if not should_move then
      return
    end
    piepan.Self:Move(e.Sender.Channel)
  end
  piepan.Audio.Play({
    filename = soundFile
  })
end)

end
