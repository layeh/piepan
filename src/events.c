/*
 * piepie - bot framework for Mumble
 *
 * Author: Tim Cooper <tim.cooper@layeh.com>
 * License: MIT (see LICENSE)
 *
 * This file contains functions that are used as event handlers for the main
 * loop.
 */

static Packet packet_read;
static Packet_Handler_Func packet_handler[26] = {
    /*  0 */ NULL,
    /*  1 */ NULL,
    /*  2 */ NULL,
    /*  3 */ NULL,
    /*  4 */ NULL,
    /*  5 */ handler_server_sync,
    /*  6 */ handler_channel_remove,
    /*  7 */ handler_channel_state,
    /*  8 */ handler_user_remove,
    /*  9 */ handler_user_state,
    /* 10 */ NULL,
    /* 11 */ handler_text_message,
    /* 12 */ NULL,
    /* 13 */ NULL,
    /* 14 */ NULL,
    /* 15 */ NULL,
    /* 16 */ NULL,
    /* 17 */ NULL,
    /* 18 */ NULL,
    /* 19 */ NULL,
    /* 20 */ NULL,
    /* 21 */ NULL,
    /* 22 */ NULL,
    /* 23 */ NULL,
    /* 24 */ handler_server_config,
    /* 25 */ NULL,
};

void
signal_event(struct ev_loop *loop, ev_signal *w, int revents)
{
    ev_break(ev_loop_main, EVBREAK_ALL);
}

void
ping_event(struct ev_loop *loop, ev_timer *w, int revents)
{
    MumbleProto__Ping ping = MUMBLE_PROTO__PING__INIT;
    sendPacket(PACKET_PING, &ping);
}

void
user_timer_event(struct ev_loop *loop, struct ev_timer *w, int revents)
{
    UserTimer *timer = (UserTimer *)w;
    lua_getglobal(timer->lua, "piepan");
    lua_getfield(timer->lua, -1, "_implOnUserTimer");
    lua_pushinteger(timer->lua, timer->id);
    lua_call(timer->lua, 1, 0);
}

void
socket_read_event(struct ev_loop *loop, ev_io *w, int revents)
{
    SSLRead *sslread = (SSLRead *)w;
    int total_read = 0;
    int ret;
    Packet_Handler_Func handler;

    ret = SSL_read(sslread->ssl, packet_read.buffer, 6);
    if (ret <= 0) {
        ev_break(loop, EVBREAK_ALL);
        return;
    }
    if (ret != 6) {
        ev_break(loop, EVBREAK_ALL);
        return;
    }
    packet_read.type = ntohs(*(uint16_t *)packet_read.buffer);
    if (packet_read.type >= sizeof(packet_handler) / sizeof(Packet_Handler_Func)) {
        ev_break(loop, EVBREAK_ALL);
        return;
    }
    packet_read.length = ntohl(*(uint32_t *)(packet_read.buffer + 2));
    if (packet_read.length > PAYLOAD_SIZE_MAX) {
        ev_break(loop, EVBREAK_ALL);
        return;
    }

    while (total_read < packet_read.length) {
        ret = SSL_read(sslread->ssl, packet_read.buffer + total_read,
                       packet_read.length - total_read);
        if (ret <= 0) {
            ev_break(loop, EVBREAK_ALL);
            return;
        }
        total_read += ret;
    }

    if (total_read != packet_read.length) {
        ev_break(loop, EVBREAK_ALL);
        return;
    }

    handler = packet_handler[packet_read.type];
    if (handler != NULL) {
        handler(sslread->lua, &packet_read);
    }
    if (SSL_pending(sslread->ssl) > 0) {
        ev_feed_fd_event(loop, w->fd, revents);
    }
}

void
audio_transmission_event(struct ev_loop *loop, struct ev_timer *w, int revents)
{
    AudioTransmission *at = (AudioTransmission *)w;
    VoicePacket packet;
    uint8_t packet_buffer[1024];
    uint8_t output[1024];
    opus_int32 byte_count;
    long long_ret;

    voicepacket_init(&packet, packet_buffer);
    voicepacket_setheader(&packet, VOICEPACKET_OPUS, VOICEPACKET_NORMAL,
        at->sequence);

    while (at->buffer.size < OPUS_FRAME_SIZE * sizeof(opus_int16)) {
        long_ret = ov_read(&at->ogg, at->buffer.pcm + at->buffer.size,
            PCM_BUFFER - at->buffer.size, 0, 2, 1, NULL);
        if (long_ret <= 0) {
            audioTransmission_stop(at, at->lua, loop);
            return;
        }
        at->buffer.size += long_ret;
    }

    byte_count = opus_encode(at->encoder, (opus_int16 *)at->buffer.pcm,
        OPUS_FRAME_SIZE, output, sizeof(output));
    if (byte_count < 0) {
        audioTransmission_stop(at, at->lua, loop);
        return;
    }
    at->buffer.size -= OPUS_FRAME_SIZE * sizeof(opus_int16);
    memmove(at->buffer.pcm, at->buffer.pcm + OPUS_FRAME_SIZE * sizeof(opus_int16),
        at->buffer.size);
    voicepacket_setframe(&packet, byte_count, output);

    sendPacketEx(PACKET_UDPTUNNEL, packet_buffer, voicepacket_getlength(&packet));

    at->sequence = (at->sequence + 1) % 10000;

    w->repeat = 0.01;
    ev_timer_again(loop, w);
}

void
script_stat_event(struct ev_loop *loop, ev_stat *w, int revents)
{
    ScriptStat *stat = (ScriptStat *)w;
    if (w->attr.st_ino == w->prev.st_ino && w->attr.st_mtime == w->prev.st_mtime) {
        return;
    }
    fprintf(stderr, "%s: reloaded %s\n", PIEPAN_NAME, stat->filename);
    lua_getglobal(stat->lua, "piepan");
    lua_getfield(stat->lua, -1, "_implLoadScript");
    lua_pushinteger(stat->lua, stat->id);
    lua_call(stat->lua, 1, 2);
    if (!lua_toboolean(stat->lua, -2)) {
        fprintf(stderr, "%s: %s\n", PIEPAN_NAME, lua_tostring(stat->lua, -1));
    }
    lua_settop(stat->lua, 0);
}
