/*
 * piepie - bot framework for Mumble
 *
 * Author: Tim Cooper <tim.cooper@layeh.com>
 * License: MIT (see LICENSE)
 *
 * This file contains functions that are used as event handlers for the main
 * loop.
 */

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
