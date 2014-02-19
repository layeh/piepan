--
-- piepie - bot framework for Mumble
--
-- Author: Tim Cooper <tim.cooper@layeh.com>
-- License: MIT (see LICENSE)
--

function piepan.internal.initialize(tbl)
    local password, tokens

    piepan.internal.state = tbl.state

    if tbl.passwordFile then
        local file, err
        if tbl.passwordFile == "-" then
            file, err = io.stdin, "could not read from stdin"
        else
            file, err = io.open(tbl.passwordFile)
        end
        if file then
            password = file:read()
            if tbl.passwordFile ~= "-" then
                file:close()
            end
        else
            print ("Error: " .. err)
        end
    end

    if tbl.tokenFile then
        local file, err = io.open(tbl.tokenFile)
        if file then
            tokens = {}
            for line in file:lines() do
                if line ~= "" then
                    table.insert(tokens, line)
                end
            end
            file:close()
        else
            print ("Error: " .. err)
        end
    end

    piepan.internal.api.apiInit(piepan.internal.api)
    piepan.internal.api.connect(tbl.username, password, tokens)
end

function piepan.internal.meta.__index(tbl, key)
    if key == "internal" then
        return
    end
    return piepan[key]
end
