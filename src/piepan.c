/*
 * piepie - bot framework for Mumble
 *
 * Author: Tim Cooper <tim.cooper@layeh.com>
 * License: MIT (see LICENSE)
 */

/* TODO:  ensure server sent a certificate and also (optionally) verify it */

#include <stdbool.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <limits.h>

#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>

#include <signal.h>
#include <netdb.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/ip.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <unistd.h>

#include <openssl/ssl.h>

#include <ev.h>

#include <opus/opus.h>
#include <vorbis/codec.h>
#include <vorbis/vorbisfile.h>

#include "../proto/Mumble.pb-c.h"

#include "piepan.h"
#include "util.c"
#include "handlers.c"
#include "events.c"
#include "api.c"
#include "piepan_impl.c"

int user_thread_pipe[2];
struct ev_loop *ev_loop_main;

static SSL_CTX *ssl_context;
static SSL *ssl;
static lua_State *lua;

int
sendPacketEx(const int type, const void *message, const int length)
{
    static Packet packet_out;
    int payload_size;
    int total_size;
    switch (type) {
        case PACKET_VERSION:
            payload_size = mumble_proto__version__get_packed_size(message);
            break;
        case PACKET_UDPTUNNEL:
            payload_size = length;
            break;
        case PACKET_AUTHENTICATE:
            payload_size = mumble_proto__authenticate__get_packed_size(message);
            break;
        case PACKET_PING:
            payload_size = mumble_proto__ping__get_packed_size(message);
            break;
        case PACKET_CHANNELREMOVE:
            payload_size = mumble_proto__channel_remove__get_packed_size(message);
            break;
        case PACKET_CHANNELSTATE:
            payload_size = mumble_proto__channel_state__get_packed_size(message);
            break;
        case PACKET_TEXTMESSAGE:
            payload_size = mumble_proto__text_message__get_packed_size(message);
            break;
        case PACKET_USERREMOVE:
            payload_size = mumble_proto__user_remove__get_packed_size(message);
            break;
        case PACKET_USERSTATE:
            payload_size = mumble_proto__user_state__get_packed_size(message);
            break;
        case PACKET_REQUESTBLOB:
            payload_size = mumble_proto__request_blob__get_packed_size(message);
            break;
        default:
            return 1;
    }
    if (payload_size >= PAYLOAD_SIZE_MAX) {
        return 2;
    }
    total_size = sizeof(uint16_t) + sizeof(uint32_t) + payload_size;
    if (payload_size > 0) {
        switch (type) {
            case PACKET_VERSION:
                mumble_proto__version__pack(message, packet_out.buffer + 6);
                break;
            case PACKET_UDPTUNNEL:
                memmove(packet_out.buffer + 6, message, length);
                break;
            case PACKET_AUTHENTICATE:
                mumble_proto__authenticate__pack(message, packet_out.buffer + 6);
                break;
            case PACKET_PING:
                mumble_proto__ping__pack(message, packet_out.buffer + 6);
                break;
            case PACKET_CHANNELREMOVE:
                mumble_proto__channel_remove__pack(message, packet_out.buffer + 6);
                break;
            case PACKET_CHANNELSTATE:
                mumble_proto__channel_state__pack(message, packet_out.buffer + 6);
                break;
            case PACKET_TEXTMESSAGE:
                mumble_proto__text_message__pack(message, packet_out.buffer + 6);
                break;
            case PACKET_USERREMOVE:
                mumble_proto__user_remove__pack(message, packet_out.buffer + 6);
                break;
            case PACKET_USERSTATE:
                mumble_proto__user_state__pack(message, packet_out.buffer + 6);
                break;
            case PACKET_REQUESTBLOB:
                mumble_proto__request_blob__pack(message, packet_out.buffer + 6);
                break;
        }
    }
    *(uint16_t *)packet_out.buffer = htons(type);
    *(uint32_t *)(packet_out.buffer + 2) = htonl(payload_size);

    return SSL_write(ssl, packet_out.buffer, total_size) == total_size ? 0 : 3;
}

void
user_thread_event(struct ev_loop *loop, ev_io *w, int revents)
{
    UserThread *user_thread;
    int thread_id;
    if (read(w->fd, &thread_id, sizeof(int)) != sizeof(int)) {
        return;
    }
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "internal");
    lua_getfield(lua, -1, "threads");
    lua_pushnumber(lua, thread_id);
    lua_gettable(lua, -2);
    lua_getfield(lua, -1, "userthread");
    user_thread = (UserThread *)lua_touserdata(lua, -1);
    if (user_thread != NULL) {
        free(user_thread);
    }
    lua_getfield(lua, -4, "events");
    lua_getfield(lua, -1, "onThreadFinish");
    lua_pushnumber(lua, thread_id);
    lua_call(lua, 1, 0);
    lua_settop(lua, 0);
}

static void
usage(FILE *stream)
{
    const char *str =
        "usage: %s [options] [scripts...]\n"
        "a bot framework for Mumble\n"
        "\n"
        "  -u <username>       username of the bot (has no effect if the certificate\n"
        "                      has been registered with the server under a different\n"
        "                      name)\n"
        "  -s <server>         address of the server (default: localhost)\n"
        "  -p <port>           port of the server (default: 64738)\n"
        "  -pw <file>          read server password from the given file (when file is -,\n"
        "                      standard input will be read)\n"
        "  -t <file>           read access tokens (one per line) from the given file\n"
        "  -c <certificate>    certificate to use for the connection\n"
        "  -k <keyfile>        key file to use for the connection (defaults to the\n"
        "                      certificate file)\n"
        "  -d                  enable development mode, which automatically reloads\n"
        "                      scripts when they are modified\n"
        "  --<name>[=<value>]  a key-value pair that will be accessible from the scripts\n"
        "  -h                  display this help\n"
        "  -v                  show version\n";
    fprintf(stream, str, PIEPAN_NAME);
}

int
main(int argc, char *argv[])
{
    struct hostent *server_host;
    char *server_host_str = "localhost";
    char *certificate_file = NULL;
    char *key_file = NULL;
    char *password_file = NULL;
    char *token_file = NULL;
    char *username = "piepan-bot";
    int port = 64738;
    int ret;
    int script_argc = -1;
    int developement_mode = 0;

    int socket_fd;
    struct sockaddr_in server_addr;

    SSLRead socket_watcher;
    ev_io user_thread_watcher;
    ev_timer ping_watcher;
    ev_signal signal_watcher;
    ev_loop_main = EV_DEFAULT;

    /*
     * Lua initialization
     */
    lua = luaL_newstate();
    if (lua == NULL) {
        fprintf(stderr, "%s: could not initialize Lua\n", PIEPAN_NAME);
        return 1;
    }
    luaL_openlibs(lua);
    if (luaL_loadbuffer(lua, (const char *)src_piepan_impl_luac,
            src_piepan_impl_luac_len, "piepan_impl") != LUA_OK) {
        fprintf(stderr, "%s: could not load piepan implementation\n", PIEPAN_NAME);
        return 1;
    }
    lua_call(lua, 0, 0);

    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "internal");
    lua_getfield(lua, -1, "api");
    lua_pushcfunction(lua, api_init);
    lua_setfield(lua, -2, "apiInit");
    lua_settop(lua, 0);

    /*
     * Argument parsing
     */
    {
        int i;
        int show_help = 0;
        int show_version = 0;
        lua_getglobal(lua, "piepan");
        lua_getfield(lua, -1, "internal");
        lua_getfield(lua, -1, "events");
        lua_getfield(lua, -1, "onArgument");
        for (i = 1; i < argc; i++) {
            int has_next = i + 1 < argc;
            if (argv[i][0] != '-' || !strcmp(argv[i], "--")) {
                script_argc = i;
                break;
            }
            if (!strcmp(argv[i], "-u") && has_next) {
                username = argv[++i];
            } else if (!strcmp(argv[i], "-c") && has_next) {
                certificate_file = argv[++i];
                if (key_file == NULL) {
                    key_file = certificate_file;
                }
            } else if (!strcmp(argv[i], "-k") && has_next) {
                key_file = argv[++i];
            } else if (!strcmp(argv[i], "-p") && has_next) {
                port = atoi(argv[++i]);
            } else if (!strcmp(argv[i], "-s") && has_next) {
                server_host_str = argv[++i];
            } else if (!strcmp(argv[i], "-h")) {
                show_help = 1;
            } else if (!strcmp(argv[i], "-pw") && has_next) {
                password_file = argv[++i];
            } else if (!strcmp(argv[i], "-t") && has_next) {
                token_file = argv[++i];
            } else if (!strcmp(argv[i], "-d")) {
                developement_mode = 1;
            } else if (!strcmp(argv[i], "-v")) {
                show_version = 1;
            } else if (!strncmp(argv[i], "--", 2) && argv[i][2] != '\0') {
                char *key = argv[i] + 2;
                char *value = strchr(key, '=');
                if (key == value) {
                    continue;
                }
                if (value != NULL) {
                    *value++ = 0;
                }
                lua_pushvalue(lua, -1);
                lua_pushstring(lua, key);
                lua_pushstring(lua, value);
                lua_call(lua, 2, 0);
            } else {
                fprintf(stderr, "%s: unknown or incomplete argument '%s'\n",
                        PIEPAN_NAME, argv[i]);
                return 1;
            }
        }
        lua_settop(lua, 0);
        if (show_version) {
            printf("%s %s (compiled on " __DATE__ " " __TIME__ ")\n", PIEPAN_NAME,
                   PIEPAN_VERSION);
            return 0;
        }
        if (show_help) {
            usage(stdout);
            return 0;
        }
    }

    /*
     * Load user scripts
     */
    {
        int i;
        lua_getglobal(lua, "piepan");
        lua_getfield(lua, -1, "internal");
        lua_getfield(lua, -1, "events");
        lua_getfield(lua, -1, "onLoadScript");
        for (i = script_argc; i >= 0 && i < argc; i++) {
            lua_pushvalue(lua, -1);
            lua_pushstring(lua, argv[i]);
            if (developement_mode) {
                lua_newuserdata(lua, sizeof(ScriptStat));
            } else {
                lua_pushnil(lua);
            }
            lua_call(lua, 2, 3);
            if (lua_toboolean(lua, -3)) {
                if (developement_mode) {
                    ScriptStat *item = lua_touserdata(lua, -1);
                    item->lua = lua;
                    item->id = lua_tointeger(lua, -2);
                    item->filename = argv[i];
                    ev_stat_init(&item->ev, script_stat_event, item->filename, 0);
                    ev_stat_start(ev_loop_main, &item->ev);
                }
            } else {
                fprintf(stderr, "%s: %s\n", PIEPAN_NAME, lua_tostring(lua, -2));
            }
            lua_pop(lua, 3);
        }
        lua_settop(lua, 0);
    }

    /*
     * Initialize Opus
     */
    {
        OpusEncoder *encoder;
        int error;

        lua_getglobal(lua, "piepan");
        lua_getfield(lua, -1, "internal");
        lua_getfield(lua, -1, "opus");
        encoder = lua_newuserdata(lua, opus_encoder_get_size(1));
        lua_setfield(lua, -2, "encoder");

        error = opus_encoder_init(encoder, 48000, 1, OPUS_APPLICATION_AUDIO);
        if (error != OPUS_OK) {
            fprintf(stderr, "%s: could not initialize the Opus encoder: %s\n",
                    PIEPAN_NAME, opus_strerror(error));
            return 1;
        }
        opus_encoder_ctl(encoder, OPUS_SET_VBR(0));
        /* TODO: set this to the server's max bitrate */
        opus_encoder_ctl(encoder, OPUS_SET_BITRATE(40000));

        lua_settop(lua, 0);
    }

    /*
     * SSL initialization
     */
    SSL_library_init();

    ssl_context = SSL_CTX_new(SSLv23_client_method());
    if (ssl_context == NULL) {
        fprintf(stderr, "%s: could not create SSL context\n", PIEPAN_NAME);
        return 1;
    }

    if (certificate_file != NULL) {
        if (!SSL_CTX_use_certificate_chain_file(ssl_context, certificate_file) ||
                !SSL_CTX_use_PrivateKey_file(ssl_context, key_file,
                                                SSL_FILETYPE_PEM) ||
                !SSL_CTX_check_private_key(ssl_context)) {
            fprintf(stderr, "%s: could not load certificate and/or key file\n",
                    PIEPAN_NAME);
            return 1;
        }
    }

    /*
     * Socket initialization and connection
     */
    socket_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (socket_fd < 0) {
        fprintf(stderr, "%s: could not create socket\n", PIEPAN_NAME);
        return 1;
    }

    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(port);

    server_host = gethostbyname(server_host_str);
    if (server_host == NULL || server_host->h_addr_list[0] == NULL ||
            server_host->h_addrtype != AF_INET) {
        fprintf(stderr, "%s: could not parse server address\n", PIEPAN_NAME);
        return 1;
    }
    memmove(&server_addr.sin_addr, server_host->h_addr_list[0],
            server_host->h_length);

    ret = connect(socket_fd, (struct sockaddr *) &server_addr,
                  sizeof(server_addr));
    if (ret != 0) {
        fprintf(stderr, "%s: could not connect to server\n", PIEPAN_NAME);
        return 1;
    }

    ssl = SSL_new(ssl_context);
    if (ssl == NULL) {
        fprintf(stderr, "%s: could not create SSL object\n", PIEPAN_NAME);
        return 1;
    }

    if (SSL_set_fd(ssl, socket_fd) == 0) {
        fprintf(stderr, "%s: could not set SSL file descriptor\n", PIEPAN_NAME);
        return 1;
    }

    if (SSL_connect(ssl) != 1) {
        fprintf(stderr, "%s: could not create secure connection\n", PIEPAN_NAME);
        return 1;
    }

    /*
     * User thread pipe
     */
    if (pipe(user_thread_pipe) != 0) {
        fprintf(stderr, "%s: could not create user thread pipe\n", PIEPAN_NAME);
        return 1;
    }

    /*
     * Trigger initial event
     */
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "internal");
    lua_getfield(lua, -1, "initialize");
    lua_newtable(lua);
    lua_pushstring(lua, username);
    lua_setfield(lua, -2, "username");
    if (password_file != NULL) {
        lua_pushstring(lua, password_file);
        lua_setfield(lua, -2, "passwordFile");
    }
    if (token_file != NULL) {
        lua_pushstring(lua, token_file);
        lua_setfield(lua, -2, "tokenFile");
    }
    lua_pushlightuserdata(lua, lua);
    lua_setfield(lua, -2, "state");
    lua_call(lua, 1, 0);
    lua_settop(lua, 0);

    /*
     * Event loop
     */
    ev_signal_init(&signal_watcher, signal_event, SIGINT);
    ev_signal_start(ev_loop_main, &signal_watcher);

    ev_io_init(&socket_watcher.ev, socket_read_event, socket_fd, EV_READ);
    socket_watcher.lua = lua;
    socket_watcher.ssl = ssl;
    ev_io_start(ev_loop_main, &socket_watcher.ev);

    ev_io_init(&user_thread_watcher, user_thread_event, user_thread_pipe[0],
               EV_READ);
    ev_io_start(ev_loop_main, &user_thread_watcher);

    ev_timer_init(&ping_watcher, ping_event, PING_TIMEOUT, PING_TIMEOUT);
    ev_timer_start(ev_loop_main, &ping_watcher);

    ev_run(ev_loop_main, 0);

    /*
     * Cleanup
     */
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "internal");
    lua_getfield(lua, -1, "events");
    lua_getfield(lua, -1, "onDisconnect");
    if (lua_isfunction(lua, -1)) {
        lua_newtable(lua);
        lua_call(lua, 1, 0);
    }

    SSL_shutdown(ssl); /* TODO:  sigpipe is triggered here if connection breaks */
    close(socket_fd);
    lua_close(lua);

    return 0;
}
