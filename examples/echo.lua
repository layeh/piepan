piepan.On('connect', function(e)
  print ("echo loaded!")
end)

piepan.On('message', function(e)
  if e.Sender == nil then
    return
  end
  piepan.Self.Channel().Send(e.Message, false)
end)
