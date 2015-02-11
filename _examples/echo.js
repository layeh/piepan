piepan.On('connect', function() {
  console.log('echo loaded!');
});

piepan.On('message', function(e) {
  if (e.Sender == null) {
    return;
  }
  piepan.Self.Channel.Send(e.Message, false);
});
