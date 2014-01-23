--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.meta.__index(table, key)
    if key == "internal" then
        return
    end
    return piepan[key]
end
