/*
 * piepie - bot framework for Mumble
 *
 * Author: Tim Cooper <tim.cooper@layeh.com>
 * License: MIT (see LICENSE)
 *
 * This file contains handlers for the messages that are received from the
 * server.
 */

void
handler_server_sync(lua_State *lua, Packet *packet)
{
    MumbleProto__ServerSync *sync =
        mumble_proto__server_sync__unpack(NULL, packet->length, packet->buffer);
    if (sync == NULL) {
        return;
    }
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "_implOnServerSync");
    if (!lua_isfunction(lua, -1)) {
        mumble_proto__server_sync__free_unpacked(sync, NULL);
        lua_settop(lua, 0);
        return;
    }
    lua_newtable(lua);
    if (sync->has_session) {
        lua_pushinteger(lua, sync->session);
        lua_setfield(lua, -2, "session");
    }
    if (sync->welcome_text != NULL) {
        lua_pushstring(lua, sync->welcome_text);
        lua_setfield(lua, -2, "welcomeText");
    }
    lua_call(lua, 1, 0);
    lua_settop(lua, 0);
    mumble_proto__server_sync__free_unpacked(sync, NULL);
}

void
handler_channel_remove(lua_State *lua, Packet *packet)
{
    MumbleProto__ChannelRemove *channel =
        mumble_proto__channel_remove__unpack(NULL, packet->length, packet->buffer);
    if (channel == NULL) {
        return;
    }
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "_implOnChannelRemove");
    if (!lua_isfunction(lua, -1)) {
        lua_settop(lua, 0);
        mumble_proto__channel_remove__free_unpacked(channel, NULL);
        return;
    }
    lua_newtable(lua);
    lua_pushinteger(lua, channel->channel_id);
    lua_setfield(lua, -2, "channelId");
    lua_call(lua, 1, 0);
    lua_settop(lua, 0);
    mumble_proto__channel_remove__free_unpacked(channel, NULL);
}

void
handler_channel_state(lua_State *lua, Packet *packet)
{
    MumbleProto__ChannelState *channel =
        mumble_proto__channel_state__unpack(NULL, packet->length, packet->buffer);
    if (channel == NULL) {
        return;
    }
    if (!channel->has_channel_id) {
        mumble_proto__channel_state__free_unpacked(channel, NULL);
        return;
    }
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "_implOnChannelState");
    if (!lua_isfunction(lua, -1)) {
        lua_settop(lua, 0);
        mumble_proto__channel_state__free_unpacked(channel, NULL);
        return;
    }
    lua_newtable(lua);
    lua_pushinteger(lua, channel->channel_id);
    lua_setfield(lua, -2, "channelId");
    if (channel->has_parent) {
        lua_pushinteger(lua, channel->parent);
        lua_setfield(lua, -2, "parentId");
    }
    if (channel->name != NULL) {
        lua_pushstring(lua, channel->name);
        lua_setfield(lua, -2, "name");
    }
    if (channel->description != NULL) {
        lua_pushstring(lua, channel->description);
        lua_setfield(lua, -2, "description");
    }
    if (channel->has_temporary) {
        lua_pushboolean(lua, channel->temporary);
        lua_setfield(lua, -2, "temporary");
    }
    lua_call(lua, 1, 0);
    lua_settop(lua, 0);
    mumble_proto__channel_state__free_unpacked(channel, NULL);
}

void
handler_server_config(lua_State *lua, Packet *packet)
{
    MumbleProto__ServerConfig *config =
       mumble_proto__server_config__unpack(NULL, packet->length, packet->buffer);
    if (config == NULL) {
        return;
    }
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "_implOnServerConfig");
    if (!lua_isfunction(lua, -1)) {
        lua_settop(lua, 0);
        mumble_proto__server_config__free_unpacked(config, NULL);
        return;
    }
    lua_newtable(lua);
    if (config->has_allow_html) {
        lua_pushboolean(lua, config->allow_html);
        lua_setfield(lua, -2, "allowHtml");
    }
    lua_call(lua, 1, 0);
    lua_settop(lua, 0);
    mumble_proto__server_config__free_unpacked(config, NULL);
}

void
handler_text_message(lua_State *lua, Packet *packet)
{
    MumbleProto__TextMessage *msg =
        mumble_proto__text_message__unpack(NULL, packet->length, packet->buffer);
    if (msg == NULL) {
        return;
    }
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "_implOnMessage");
    if (!lua_isfunction(lua, -1)) {
        lua_settop(lua, 0);
        mumble_proto__text_message__free_unpacked(msg, NULL);
        return;
    }
    lua_newtable(lua);
    if (msg->has_actor) {
        lua_pushinteger(lua, msg->actor);
        lua_setfield(lua, -2, "actor");
    }
    if (msg->message != NULL) {
        lua_pushstring(lua, msg->message);
        lua_setfield(lua, -2, "message");
    }
    if (msg->n_session > 0) {
        int i;
        lua_newtable(lua);
        for (i = 0; i < msg->n_session; i++) {
            lua_pushinteger(lua, i);
            lua_pushinteger(lua, msg->session[i]);
            lua_settable(lua, -3);
        }
        lua_setfield(lua, -2, "users");
    }
    if (msg->n_channel_id > 0) {
        int i;
        lua_newtable(lua);
        for (i = 0; i < msg->n_channel_id; i++) {
            lua_pushinteger(lua, i);
            lua_pushinteger(lua, msg->channel_id[i]);
            lua_settable(lua, -3);
        }
        lua_setfield(lua, -2, "channels");
    }
    lua_call(lua, 1, 0);
    lua_settop(lua, 0);
    mumble_proto__text_message__free_unpacked(msg, NULL);
}

void
handler_user_state(lua_State *lua, Packet *packet)
{
    MumbleProto__UserState *user =
        mumble_proto__user_state__unpack(NULL, packet->length, packet->buffer);
    if (user == NULL) {
        return;
    }
    if (!user->has_session) {
        mumble_proto__user_state__free_unpacked(user, NULL);
        return;
    }

    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "_implOnUserChange");
    if (!lua_isfunction(lua, -1)) {
        lua_settop(lua, 0);
        mumble_proto__user_state__free_unpacked(user, NULL);
        return;
    }
    lua_newtable(lua);
    lua_pushinteger(lua, user->session);
    lua_setfield(lua, -2, "session");
    if (user->has_actor) {
        lua_pushinteger(lua, user->actor);
        lua_setfield(lua, -2, "actor");
    }
    if (user->name != NULL) {
        lua_pushstring(lua, user->name);
        lua_setfield(lua, -2, "name");
    }
    if (user->has_channel_id) {
        lua_pushinteger(lua, user->channel_id);
        lua_setfield(lua, -2, "channelId");
    }
    if (user->has_user_id) {
        lua_pushinteger(lua, user->user_id);
        lua_setfield(lua, -2, "userId");
    }
    if (user->has_mute) {
        lua_pushboolean(lua, user->mute);
        lua_setfield(lua, -2, "isServerMuted");
    }
    if (user->has_deaf) {
        lua_pushboolean(lua, user->deaf);
        lua_setfield(lua, -2, "isServerDeafened");
    }
    if (user->has_self_mute) {
        lua_pushboolean(lua, user->self_mute);
        lua_setfield(lua, -2, "isSelfMuted");
    }
    if (user->has_self_deaf) {
        lua_pushboolean(lua, user->self_deaf);
        lua_setfield(lua, -2, "isSelfDeafened");
    }
    if (user->has_suppress) {
        lua_pushboolean(lua, user->suppress);
        lua_setfield(lua, -2, "isSuppressed");
    }
    if (user->comment != NULL) {
        lua_pushstring(lua, user->comment);
        lua_setfield(lua, -2, "comment");
    }
    if (user->has_recording) {
        lua_pushboolean(lua, user->recording);
        lua_setfield(lua, -2, "isRecording");
    }
    lua_call(lua, 1, 0);
    lua_settop(lua, 0);
    mumble_proto__user_state__free_unpacked(user, NULL);
}

void
handler_user_remove(lua_State *lua, Packet *packet)
{
    MumbleProto__UserRemove *user =
       mumble_proto__user_remove__unpack(NULL, packet->length, packet->buffer);
    if (user == NULL) {
        return;
    }
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "_implOnUserRemove");
    if (!lua_isfunction(lua, -1)) {
        lua_settop(lua, 0);
        mumble_proto__user_remove__free_unpacked(user, NULL);
        return;
    }
    lua_newtable(lua);
    lua_pushinteger(lua, user->session);
    lua_setfield(lua, -2, "session");
    if (user->has_actor) {
        lua_pushinteger(lua, user->actor);
        lua_setfield(lua, -2, "actor");
    }
    if (user->reason != NULL) {
        lua_pushstring(lua, user->reason);
        lua_setfield(lua, -2, "reason");
    }
    if (user->has_ban) {
        lua_pushboolean(lua, user->ban);
        lua_setfield(lua, -2, "ban");
    }
    lua_call(lua, 1, 0);
    lua_settop(lua, 0);
    mumble_proto__user_remove__free_unpacked(user, NULL);
}
