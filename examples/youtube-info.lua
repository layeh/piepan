--
-- Echos back any links to youtube videos with the video's title, duration, and
-- thumbnail
--
-- Requires:  wget, jshon
--

function piepan.onConnect()
    print ("youtube-info loaded!")
end

function piepan.onMessage(message)
    if message.user == nil then
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
        local video_id = string.match(message.text, pattern)
        if video_id ~= nil and string.len(video_id) < 20 then
            piepan.Thread.new(youtube_info_lookup, youtube_info_completed,
                video_id)
            return
        end
    end
end

function youtube_info_lookup(id)
    if id == nil then
        return
    end
    local cmd = [[
        wget -q -O - 'http://gdata.youtube.com/feeds/api/videos/%s?v=2&alt=jsonc' |
            jshon -Q -e data -e title -u -p -e duration -u -p -e thumbnail -e hqDefault -u
    ]]
    local jshon = io.popen(string.format(cmd, id))
    local name = jshon:read()
    local duration = jshon:read()
    local thumbnail = jshon:read()
    if name == nil or duration == nil then
        return
    end

    return {
        id = id,
        title = name,
        duration = string.format("%d:%02d", duration / 60, duration % 60),
        thumbnail = thumbnail
    }
end

function youtube_info_completed(info)
    if info == nil then
        return
    end
    local fmt = [[
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
    local message = string.format(fmt, info.id, info.title, info.duration,
        info.id, info.thumbnail)
    piepan.me.channel:send(message)
end
