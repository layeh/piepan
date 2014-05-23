/*
 * piepie - bot framework for Mumble
 *
 * Author: Tim Cooper <tim.cooper@layeh.com>
 * License: MIT (see LICENSE)
 *
 */

#define PIEPAN_NAME "piepan"
#define PIEPAN_VERSION "0.1.1"

#define PAYLOAD_SIZE_MAX (1024 * 1024 * 8 - 1)

#ifndef PING_TIMEOUT
#define PING_TIMEOUT 15.0
#endif

#define OPUS_FRAME_SIZE 480
#define PCM_BUFFER 4096

enum {
    PACKET_VERSION          = 0,
    PACKET_UDPTUNNEL        = 1,
    PACKET_AUTHENTICATE     = 2,
    PACKET_PING             = 3,
    PACKET_CHANNELREMOVE    = 6,
    PACKET_CHANNELSTATE     = 7,
    PACKET_USERREMOVE       = 8,
    PACKET_USERSTATE        = 9,
    PACKET_TEXTMESSAGE      = 11,
    PACKET_REQUESTBLOB      = 23
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
    ev_io ev;
    lua_State *lua;
    SSL *ssl;
} SSLRead;

typedef struct {
    ev_timer ev;
    lua_State *lua;
    int id;
} UserTimer;

typedef struct {
    lua_State *lua;
    pthread_t thread;
    int id;
} UserThread;

typedef struct {
    ev_timer ev;
    FILE *file;
    lua_State *lua;
    OggVorbis_File ogg;
    uint32_t sequence;
    OpusEncoder *encoder;
    struct {
        char pcm[PCM_BUFFER];
        int size;
    } buffer;
} AudioTransmission;

typedef struct {
    uint8_t *buffer;
    int length;
    int header_length;
} VoicePacket;

typedef struct {
    ev_stat ev;
    lua_State *lua;
    int id;
    char *filename;
} ScriptStat;

/*
 * Prototypes
 */
#define sendPacket(type, message) sendPacketEx(type, message, 0)
int sendPacketEx(const int type, const void *message, const int length);

typedef void (*Packet_Handler_Func)(lua_State *lua, Packet *packet);
void user_timer_event(struct ev_loop *loop, struct ev_timer *w, int revents);
void audio_transmission_event(struct ev_loop *loop, struct ev_timer *w, int revents);
// TODO:  remove globals -- pass important global data to Lua
extern struct ev_loop *ev_loop_main;
extern int user_thread_pipe[2];
