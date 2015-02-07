/*
 * Soundboard example.
 *
 * Sounds are triggered when #<keyword> appears in a text message, where
 * keyword has been defined below in the sounds table.
 */

(function() {
  // Boolean if users need to be registered on the server to trigger sounds
  var require_registered = true;

  // Boolean if sounds should stop playing when another is triggered
  var interrupt_sounds = false;

  // Boolean if the bot should move into the user's channel to play the sound
  var should_move = false;

  // Table with keys being the keywords and values being the sound files
  var sounds = {
    cheer: "sounds/cheer.ogg",
    hello: "sounds/hello.ogg",
    huh:   "sounds/huh.ogg",
    image: "sounds/image.ogg",
    lol:   "sounds/lol.ogg",
    mock:  "sounds/mock.ogg",
    nice:  "sounds/nice.ogg"
  };

  // Sound file path prefix
  var prefix = ENV['PREFIX'] || "_examples/";

  piepan.On('connect', function() {
    console.log("Soundboard loaded!");
  });

  piepan.On('message', function(e) {
    if (e.Sender == null) {
      return;
    }

    var search = e.Message.match(/#(\w+)/);
    if (!search || !sounds[search[1]]) {
      return;
    }
    var soundFile = prefix + sounds[search[1]];
    if (require_registered && !e.Sender.IsRegistered()) {
      e.Sender.Send("You must be registered on the server to trigger sounds.");
      return;
    }
    if (piepan.Audio.IsPlaying() && !interrupt_sounds) {
      return;
    }
    if (piepan.Self.Channel().ID() != e.Sender.Channel().ID()) {
      if (!should_move) {
        return;
      }
      piepan.Self.Move(e.Sender.Channel());
    }

    piepan.Audio.Stop();
    piepan.Audio.Play({
      filename: soundFile,
    });
  });
})();
