printf = function(s,...)
    return io.write(s:format(...))
end
function file_exists(name)
   local f=io.open(name,"r")
   if f~=nil then io.close(f) return true else return false end
end

playNext = function(file_name, e, idx)
    printf("beep boop: %i\n", idx)
    if e.Client.Channels[idx] == nil then
        print("done")
        os.exit()
    else
        print("moving into channel")
        piepan.Self:Move(e.Client.Channels[idx])
        piepan.Audio.New({filename = file_name, callback = function()
            playNext(file_name, e, idx + 1)
        end}):Play()
    end
end
if piepan.Args == nil or not #piepan.Args == 1 then
    print("Requires a script argument for sound to play")
    os.exit()
end
file_name = piepan.Args[1]
if file_exists(file_name) then
    piepan.On ('connect', function(e)
        playNext(file_name, e, 0)
    end)
else
    print("Unable to open file")
end
