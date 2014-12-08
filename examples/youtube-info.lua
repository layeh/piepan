--
-- Echos back any links to youtube videos with the video's title, duration, and
-- thumbnail
--
-- Requires:  wget, jshon
--

local prefix = os.getenv('PREFIX') or "examples/"
local worker = prefix .. "youtube-info-worker.sh"

local message_fmt = [[
<table>
    <tr>
        <td valign="middle">
            <img src='https://www.youtube.com/yt/brand/media/image/YouTube-icon-full_color.png' height="25" />
        </td>
        <td align="center" valign="middle">
            <a href="http://youtu.be/%s">%s (%s)</a>
        </td>
    </tr>
    <tr>
        <td></td>
        <td align="center">
            <a href="http://youtu.be/%s"><img src="%s" width="250" /></a>
        </td>
    </tr>
</table>
]]

piepan.On('connect', function()
  print ("youtube-info loaded!")
end)

piepan.On('message', function(e)
  if e.Sender == nil then
    return
  end
  local patterns = {
    "https?://www%.youtube%.com/watch%?v=([%d%a_%-]+)",
    "https?://youtube%.com/watch%?v=([%d%a_%-]+)",
    "https?://youtu.be/([%d%a_%-]+)",
    "https?://youtube.com/v/([%d%a_%-]+)",
    "https?://www.youtube.com/v/([%d%a_%-]+)"
  }
  for _,pattern in ipairs(patterns) do
    local video_id = string.match(e.Message, pattern)
    if video_id ~= nil and string.len(video_id) < 20 then
      piepan.Process.New(function (success, data)
        if not success then
          return
        end

        local id, title, duration, thumbnail = string.match(data, "([^\r\n]+)\r?\n([^\r\n]+)\r?\n([^\r\n]+)\r?\n([^\r\n]+)\r?\n")
        duration = string.format("%d:%02d", duration / 60, duration % 60)
        local message = string.format(message_fmt, id, title, duration,
          id, thumbnail)
        piepan.Self.Channel().Send(message, false)
      end, worker, video_id)
      return
    end
  end
end)
