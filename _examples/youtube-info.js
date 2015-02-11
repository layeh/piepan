/*
 * Echos back any links to youtube videos with the video's title, duration,
 * and thumbnail.
 *
 * Requires:  wget
 */

(function() {

var message_template = _.template('\
<table>\
    <tr>\
        <td valign="middle">\
            <img src="https://www.youtube.com/yt/brand/media/image/YouTube-icon-full_color.png" height="25" />\
        </td>\
        <td align="center" valign="middle">\
            <a href="https://youtu.be/<%= id %>"><%= title %> (<%= duration %>)</a>\
        </td>\
    </tr>\
    <tr>\
        <td></td>\
        <td align="center">\
            <a href="https://youtu.be/<%= id %>"><img src="<%= thumbnail %>" width="250" /></a>\
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
      var json = JSON.parse(data);
      var seconds = json.data.duration;
      var minutes = Math.floor(seconds / 60).toFixed(0).toString();
      seconds = (seconds % 60).toFixed(0);
      var duration = minutes + ":";
      if (seconds < 10) {
        duration += "0";
      }
      duration += seconds;

      var message = message_template({
        id: json.data.id,
        title: json.data.title,
        duration: duration,
        thumbnail: json.data.thumbnail.hqDefault,
      });
      piepan.Self.Channel.Send(message, false);
    }, 'wget', '-q', '-O', '-', 'http://gdata.youtube.com/feeds/api/videos/' + video_id + '?v=2&alt=jsonc');
    return;
  }
});

})();
