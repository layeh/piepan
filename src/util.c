/*
 * piepie - bot framework for Mumble
 *
 * Author: Tim Cooper <tim.cooper@layeh.com>
 * License: MIT (see LICENSE)
 *
 * This file contains utility functions.
 */

int
util_set_varint(uint8_t buffer[], const uint64_t value)
{
    if (value < 0x80) {
        buffer[0] = value;
        return 1;
    } else if (value < 0x4000) {
        buffer[0] = (value >> 8) | 0x80;
        buffer[1] = value & 0xFF;
        return 2;
    }
    return -1;
}

void
audioTransmission_stop(AudioTransmission *at, lua_State *lua, struct ev_loop *loop)
{
    if (at == NULL || lua == NULL || loop == NULL) {
        return;
    }
    ov_clear(&at->ogg);
    fclose(at->file);
    ev_timer_stop(loop, &at->ev);

    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "internal");
    lua_getfield(lua, -1, "events");
    lua_getfield(lua, -1, "onAudioFinished");
    lua_call(lua, 0, 0);
    lua_settop(lua, 0);
}

#define VOICEPACKET_NORMAL 0
#define VOICEPACKET_OPUS 4

VoicePacket *
voicepacket_init(VoicePacket *packet, uint8_t *buffer)
{
    if (packet == NULL || buffer == NULL) {
        return NULL;
    }
    packet->buffer = buffer;
    packet->length = 0;
    packet->header_length = 0;
    return packet;
}

int
voicepacket_setheader(VoicePacket *packet, const uint8_t type,
        const uint8_t target, const uint32_t sequence)
{
    int offset;
    if (packet == NULL) {
        return -1;
    }
    if (packet->length > 0) {
        return -2;
    }
    packet->buffer[0] = ((type & 0x7) << 5) | (target & 0x1F);
    offset = util_set_varint(packet->buffer + 1, sequence);
    packet->length = packet->header_length = 1 + offset;
    return 1;
}

int
voicepacket_setframe(VoicePacket *packet, const uint16_t length, uint8_t *buffer)
{
    int offset;
    if (packet == NULL || buffer == NULL || length <= 0 || length >= 0x2000) {
        return -1;
    }
    if (packet->header_length <= 0) {
        return -2;
    }
    offset = util_set_varint(packet->buffer + packet->header_length, length);
    if (offset <= 0) {
        return -3;
    }
    memmove(packet->buffer + packet->header_length + offset, buffer, length);
    packet->length = packet->header_length + length + offset;
    return 1;
}

int
voicepacket_getlength(const VoicePacket *packet)
{
    if (packet == NULL) {
        return -1;
    }
    return packet->length;
}
