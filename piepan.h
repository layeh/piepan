/*
 * piepie - bot framework for Mumble
 *
 * Author: Tim Cooper <tim.cooper@layeh.com>
 * License: MIT (see LICENSE)
 *
 */

#define PAYLOAD_SIZE_MAX (1024 * 1024 * 8 - 1)
#define PIEPAN_VERSION "0.1.1"

enum {
    PACKET_VERSION      = 0,
    PACKET_AUTHENTICATE = 2,
    PACKET_PING         = 3,
    PACKET_USERREMOVE   = 8,
    PACKET_USERSTATE    = 9,
    PACKET_TEXTMESSAGE  = 11
};

/*
 * Structures
 */

typedef struct {
    uint16_t type;
    uint32_t length;
    uint8_t buffer[PAYLOAD_SIZE_MAX + 6];
} Packet;

typedef struct {
    ev_timer ev;
    int id;
} UserTimer;

typedef struct {
    lua_State *lua;
    pthread_t thread;
    int id;
} UserThread;

/*
 * Prototypes
 */
int sendPacket(int type, void *message);

typedef void (*Packet_Handler_Func)(lua_State *lua, Packet *packet);
void user_timer_event(struct ev_loop *loop, struct ev_timer *w, int revents);
extern struct ev_loop *ev_loop_main;
extern int user_thread_pipe[2];
