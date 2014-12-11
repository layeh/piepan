/*
 * Echos back any links to youtube videos with the video's title, duration,
 * and thumbnail.
 *
 * Requires:  wget, jshon
 */

(function() {

var prefix = ENV['PREFIX'] || 'examples/';
var worker = prefix + 'youtube-info-worker.sh';

var message_template = _.template('\
<table>\
    <tr>\
        <td valign="middle">\
            <img src="https://www.youtube.com/yt/brand/media/image/YouTube-icon-full_color.png" height="25" />\
        </td>\
        <td align="center" valign="middle">\
            <a href="http://youtu.be/<%= id %>"><%= title %> (<%= duration %>)</a>\
        </td>\
    </tr>\
    <tr>\
        <td></td>\
        <td align="center">\
            <a href="http://youtu.be/<%= id %>"><img src="<%= thumbnail %>" width="250" /></a>\
        </td>\
    </tr>\
</table>');

piepan.On('connect', function() {
  console.log('youtube-info loaded!')
});

piepan.On('message', function(e) {
  if (e.Sender == null) {
    return;
  }
  var patterns = [
    /https?:\/\/www\.youtube\.com\/watch\?v=([\w-]+)/,
    /https?:\/\/youtube\.com\/watch\?v=([\w-]+)/,
    /https?:\/\/youtu.be\/([\w-]+)/,
    /https?:\/\/youtube.com\/v\/([\w-]+)/,
    /https?:\/\/www.youtube.com\/v\/([\w-]+)/
  ];
  for (var i = 0; i < patterns.length; i++) {
    var pattern = patterns[i];
    var matches = e.Message.match(pattern);
    if (!matches) {
      continue;
    }
    var video_id = matches[1];
    if (video_id.length >= 20) {
      continue;
    }

    piepan.Process.New(function (success, data) {
      if (!success) {
        return;
      }

      var matches = data.match(/([^\r\n]+)\r?\n([^\r\n]+)\r?\n([^\r\n]+)\r?\n([^\r\n]+)\r?\n/);

      var minutes = (matches[3] / 60).toFixed(0).toString();
      var seconds = (matches[3] % 60).toFixed(0);
      var duration = minutes + ":";
      if (seconds < 10) {
        duration += "0";
      }
      duration += seconds;

      var message = message_template({
        id: matches[1],
        title: matches[2],
        duration: duration,
        thumbnail: matches[4],
      });
      piepan.Self.Channel().Send(message, false);
    }, worker, video_id);
    return;
  }
});

})();
