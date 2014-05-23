/*
 * piepie - bot framework for Mumble
 *
 * Author: Tim Cooper <tim.cooper@layeh.com>
 * License: MIT (see LICENSE)
 *
 * This file contains native functions that are called from the piepan Lua
 * script.
 */

// TODO:  get rid of selfs and only pass non-tables?
// TODO:  use Lua user data in place of mallocing it ourselves?

#include <pthread.h>

int
api_User_send(lua_State *lua)
{
    // [self, message]
    MumbleProto__TextMessage msg = MUMBLE_PROTO__TEXT_MESSAGE__INIT;
    uint32_t session;
    lua_getfield(lua, -2, "session");
    session = lua_tointeger(lua, -1);
    msg.n_session = 1;
    msg.session = &session;
    msg.message = (char *)lua_tostring(lua, -2);
    sendPacket(PACKET_TEXTMESSAGE, &msg);
    return 0;
}

int
api_User_kick(lua_State *lua)
{
    // [self, string reason]
    MumbleProto__UserRemove msg = MUMBLE_PROTO__USER_REMOVE__INIT;
    lua_getfield(lua, -2, "session");
    msg.session = lua_tointeger(lua, -1);
    if (lua_isstring(lua, -2)) {
        msg.reason = (char *)lua_tostring(lua, -2);
    }
    sendPacket(PACKET_USERREMOVE, &msg);
    return 0;
}

int
api_User_ban(lua_State *lua)
{
    // [self, string reason]
    MumbleProto__UserRemove msg = MUMBLE_PROTO__USER_REMOVE__INIT;
    lua_getfield(lua, -2, "session");
    msg.session = lua_tointeger(lua, -1);
    if (lua_isstring(lua, -2)) {
        msg.reason = (char *)lua_tostring(lua, -2);
    }
    msg.has_ban = true;
    msg.ban = true;
    sendPacket(PACKET_USERREMOVE, &msg);
    return 0;
}

int
api_User_moveTo(lua_State *lua)
{
    // [self, int channel_id]
    MumbleProto__UserState msg = MUMBLE_PROTO__USER_STATE__INIT;
    msg.channel_id = lua_tointeger(lua, -1);
    lua_getfield(lua, -2, "session");
    msg.session = lua_tointeger(lua, -1);
    sendPacket(PACKET_USERSTATE, &msg);
    return 0;
}

int
api_User_setComment(lua_State *lua)
{
    // [self, string comment]
    MumbleProto__UserState msg = MUMBLE_PROTO__USER_STATE__INIT;
    msg.comment = (char *)lua_tostring(lua, -1);
    lua_getfield(lua, -2, "session");
    msg.session = lua_tointeger(lua, -1);
    sendPacket(PACKET_USERSTATE, &msg);
    return 0;
}

int
api_User_register(lua_State *lua)
{
    // [self]
    MumbleProto__UserState msg = MUMBLE_PROTO__USER_STATE__INIT;
    msg.has_user_id = true;
    msg.user_id = 0;
    sendPacket(PACKET_USERSTATE, &msg);
    return 0;
}

int
api_User_setTexture(lua_State *lua)
{
    // [self, texture]
    MumbleProto__UserState msg = MUMBLE_PROTO__USER_STATE__INIT;
    msg.has_texture = true;
    msg.texture.data = (uint8_t *)lua_tolstring(lua, -1, &msg.texture.len);
    lua_getfield(lua, 1, "session");
    msg.session = lua_tointeger(lua, -1);
    sendPacket(PACKET_USERSTATE, &msg);
    return 0;
}

int
api_Channel_play(lua_State *lua)
{
    // [lua_State *, OpusEncoder *encoder, string filename]
    AudioTransmission *at = malloc(sizeof(AudioTransmission));
    if (at == NULL) {
        return 0;
    }
    at->file = fopen(lua_tostring(lua, 3), "rb");
    if (at->file == NULL) {
        free(at);
        return 0;
    }
    if (ov_open_callbacks(at->file, &at->ogg, NULL, 0, OV_CALLBACKS_STREAMONLY_NOCLOSE) != 0) {
        fclose(at->file);
        free(at);
        return 0;
    }

    at->lua = lua_touserdata(lua, 1);
    at->encoder = lua_touserdata(lua, 2);
    at->sequence = 1;
    at->buffer.size = 0;
    ev_timer_init(&at->ev, audio_transmission_event, 0., 0.01);
    ev_timer_start(ev_loop_main, &at->ev);

    lua_pushlightuserdata(lua, at);
    return 1;
}

int
api_Channel_send(lua_State *lua)
{
    // [self, message]
    MumbleProto__TextMessage msg = MUMBLE_PROTO__TEXT_MESSAGE__INIT;
    uint32_t channel;
    lua_getfield(lua, 1, "id");
    msg.message = (char *)lua_tostring(lua, -2);
    channel = lua_tointeger(lua, -1);
    msg.n_channel_id = 1;
    msg.channel_id = &channel;
    sendPacket(PACKET_TEXTMESSAGE, &msg);
    return 0;
}

int
api_Channel_setDescription(lua_State *lua)
{
    // [self, description]
    MumbleProto__ChannelState msg = MUMBLE_PROTO__CHANNEL_STATE__INIT;
    msg.has_channel_id = true;
    lua_getfield(lua, 1, "id");
    msg.channel_id = lua_tointeger(lua, -1);
    msg.description = (char *)lua_tostring(lua, 2);
    sendPacket(PACKET_CHANNELSTATE, &msg);
    return 0;
}

int
api_Channel_remove(lua_State *lua)
{
    // [self]
    MumbleProto__ChannelRemove msg = MUMBLE_PROTO__CHANNEL_REMOVE__INIT;
    lua_getfield(lua, 1, "id");
    msg.channel_id = lua_tointeger(lua, -1);
    sendPacket(PACKET_CHANNELREMOVE, &msg);
    return 0;
}

int
api_Timer_new(lua_State *lua)
{
    // [table, id, timeout, lua_State]
    UserTimer *timer;
    timer = lua_newuserdata(lua, sizeof(UserTimer));
    lua_setfield(lua, 1, "ptr");
    timer->id = lua_tonumber(lua, 2);
    timer->lua = lua_touserdata(lua, 4);

    ev_timer_init(&timer->ev, user_timer_event, lua_tonumber(lua, 3), 0.);
    ev_timer_start(ev_loop_main, &timer->ev);

    return 0;
}

int
api_Timer_cancel(lua_State *lua)
{
    // [UserTimer *]
    UserTimer *timer = lua_touserdata(lua, -1);
    ev_timer_stop(ev_loop_main, &timer->ev);
    return 0;
}

static void *
api_Thread_worker(void *arg)
{
    UserThread *user_thread = (UserThread *)arg;
    lua_getglobal(user_thread->lua, "piepan");
    lua_getfield(user_thread->lua, -1, "internal");
    lua_getfield(user_thread->lua, -1, "events");
    lua_getfield(user_thread->lua, -1, "onThreadExecute");
    lua_pushnumber(user_thread->lua, user_thread->id);
    lua_call(user_thread->lua, 1, 0);
    write(user_thread_pipe[1], &user_thread->id, sizeof(int));
    return NULL;
}

int
api_Thread_new(lua_State *lua)
{
    // [Thread, id]
    UserThread *user_thread;
    user_thread = (UserThread *)malloc(sizeof(UserThread));
    if (user_thread == NULL) {
        return 0;
    }
    user_thread->id = lua_tonumber(lua, -1);
    user_thread->lua = lua_newthread(lua);
    lua_setfield(lua, -3, "lua");
    lua_pushlightuserdata(lua, user_thread);
    lua_setfield(lua, -3, "userthread");
    pthread_create(&user_thread->thread, NULL, api_Thread_worker, user_thread);
    return 0;
}

int
api_Audio_stop(lua_State *lua)
{
    // [AudioTransmission *]
    AudioTransmission *at = (AudioTransmission *)lua_touserdata(lua, -1);
    audioTransmission_stop(at, lua, ev_loop_main);
    return 0;
}

int
api_disconnect(lua_State *lua)
{
    kill(0, SIGINT);
    return 0;
}

int
api_connect(lua_State *lua)
{
    // [string username, string password, table tokens]
    MumbleProto__Version version = MUMBLE_PROTO__VERSION__INIT;
    MumbleProto__Authenticate auth = MUMBLE_PROTO__AUTHENTICATE__INIT;

    auth.has_opus = true;
    auth.opus = true;
    auth.username = (char *)lua_tostring(lua, -3);

    if (!lua_isnil(lua, -2)) {
        auth.password = (char *)lua_tostring(lua, -2);
    }
    if (!lua_isnil(lua, -1)) {
        lua_len(lua, -1);
        auth.n_tokens = lua_tointeger(lua, -1);
        lua_pop(lua, 1);

        if (lua_checkstack(lua, auth.n_tokens)) {
            int i;
            int table_index = lua_absindex(lua, -1);
            auth.tokens = lua_newuserdata(lua, sizeof(char *) * auth.n_tokens);
            lua_pushnil(lua);

            for (i = 0; i < auth.n_tokens; i++) {
                if (!lua_next(lua, table_index)) {
                    break;
                }
                auth.tokens[i] = (char *)lua_tostring(lua, -1);
                lua_insert(lua, -2);
            }

            auth.n_tokens = i;
        } else {
            // TODO:  notify the user of this
            auth.n_tokens = 0;
        }
    }
    version.has_version = true;
    version.version = 1 << 16 | 2 << 8 | 5; // 1.2.5
    version.release = "Unknown";
    version.os = PIEPAN_NAME;
    version.os_version = PIEPAN_VERSION;

    sendPacket(PACKET_VERSION, &version);
    sendPacket(PACKET_AUTHENTICATE, &auth);
    return 0;
}

int
api_resolveHashes(lua_State *lua)
{
    // [sessionTextures, sessionComments, channelDescriptions]
    MumbleProto__RequestBlob msg = MUMBLE_PROTO__REQUEST_BLOB__INIT;
    int i;

    if (!lua_isnil(lua, 1) && luaL_len(lua, 1) > 0) {
        msg.n_session_texture = luaL_len(lua, 1);
        msg.session_texture = lua_newuserdata(lua,
            sizeof(uint32_t) * msg.n_session_texture);
        lua_pushnil(lua);
        for (i = 0; i < msg.n_session_texture; i++) {
            if (!lua_next(lua, 1)) {
                break;
            }
            msg.session_texture[i] = lua_tonumber(lua, -1);
            lua_pop(lua, 1);
        }
        lua_pop(lua, 1);
    }

    if (!lua_isnil(lua, 2) && luaL_len(lua, 2) > 0) {
        msg.n_session_comment = luaL_len(lua, 2);
        msg.session_comment = lua_newuserdata(lua,
            sizeof(uint32_t) * msg.n_session_comment);
        lua_pushnil(lua);
        for (i = 0; i < msg.n_session_comment; i++) {
            if (!lua_next(lua, 2)) {
                break;
            }
            msg.session_comment[i] = lua_tonumber(lua, -1);
            lua_pop(lua, 1);
        }
        lua_pop(lua, 1);
    }

    if (!lua_isnil(lua, 3) && luaL_len(lua, 3) > 0) {
        msg.n_channel_description = luaL_len(lua, 3);
        msg.channel_description = lua_newuserdata(lua,
            sizeof(uint32_t) * msg.n_channel_description);
        lua_pushnil(lua);
        for (i = 0; i < msg.n_channel_description; i++) {
            if (!lua_next(lua, 3)) {
                break;
            }
            msg.channel_description[i] = lua_tonumber(lua, -1);
            lua_pop(lua, 1);
        }
        lua_pop(lua, 1);
    }

    sendPacket(PACKET_REQUESTBLOB, &msg);
    return 0;
}

int
api_init(lua_State *lua)
{
    // [table]

    lua_pushcfunction(lua, api_User_send);
    lua_setfield(lua, -2, "userSend");
    lua_pushcfunction(lua, api_User_kick);
    lua_setfield(lua, -2, "userKick");
    lua_pushcfunction(lua, api_User_ban);
    lua_setfield(lua, -2, "userBan");
    lua_pushcfunction(lua, api_User_moveTo);
    lua_setfield(lua, -2, "userMoveTo");
    lua_pushcfunction(lua, api_User_setComment);
    lua_setfield(lua, -2, "userSetComment");
    lua_pushcfunction(lua, api_User_setTexture);
    lua_setfield(lua, -2, "userSetTexture");
    lua_pushcfunction(lua, api_User_register);
    lua_setfield(lua, -2, "userRegister");

    lua_pushcfunction(lua, api_Channel_play);
    lua_setfield(lua, -2, "channelPlay");
    lua_pushcfunction(lua, api_Channel_send);
    lua_setfield(lua, -2, "channelSend");
    lua_pushcfunction(lua, api_Channel_setDescription);
    lua_setfield(lua, -2, "channelSetDescription");
    lua_pushcfunction(lua, api_Channel_remove);
    lua_setfield(lua, -2, "channelRemove");

    lua_pushcfunction(lua, api_Timer_new);
    lua_setfield(lua, -2, "timerNew");
    lua_pushcfunction(lua, api_Timer_cancel);
    lua_setfield(lua, -2, "timerCancel");

    lua_pushcfunction(lua, api_Thread_new);
    lua_setfield(lua, -2, "threadNew");

    lua_pushcfunction(lua, api_Audio_stop);
    lua_setfield(lua, -2, "audioStop");

    lua_pushcfunction(lua, api_connect);
    lua_setfield(lua, -2, "connect");

    lua_pushcfunction(lua, api_resolveHashes);
    lua_setfield(lua, -2, "resolveHashes");

    lua_pushcfunction(lua, api_disconnect);
    lua_setfield(lua, -2, "disconnect");
    return 0;
}
