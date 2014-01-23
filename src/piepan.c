/*
 * piepie - bot framework for Mumble
 *
 * Author: Tim Cooper <tim.cooper@layeh.com>
 * License: MIT (see LICENSE)
 */

// TODO:  ensure server sent a certificate and also (optionally) verify it

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


#ifndef PING_TIMEOUT
#define PING_TIMEOUT 15.0
#endif

#ifndef TOKEN_BUFFER_SIZE
#define TOKEN_BUFFER_SIZE 1024
#endif

#ifndef MAX_TOKENS
#define MAX_TOKENS 32
#endif

int user_thread_pipe[2];
struct ev_loop *ev_loop_main;

static const char *progname;
static SSL_CTX *ssl_context;
static SSL *ssl;
static lua_State *lua;
static Packet packet_out;

typedef struct ScriptStat {
    ev_stat ev;
    int id;
    char *filename;
    struct ScriptStat *next;
} ScriptStat;

static const char *
impl_reader(lua_State *L, void *data, size_t *size)
{
    int *ret = (int *)data;
    if (*ret == 0) {
        *ret = 1;
        *size = src_piepan_impl_luac_len;
        return (const char *)src_piepan_impl_luac;
    } else {
        return NULL;
    }
}

int
sendPacketEx(int type, void *data, int length)
{
    int payload_size;
    int total_size;
    switch (type) {
        case PACKET_VERSION:
            payload_size = mumble_proto__version__get_packed_size(data);
            break;
        case PACKET_UDPTUNNEL:
            payload_size = length;
            break;
        case PACKET_AUTHENTICATE:
            payload_size = mumble_proto__authenticate__get_packed_size(data);
            break;
        case PACKET_PING:
            payload_size = mumble_proto__ping__get_packed_size(data);
            break;
        case PACKET_TEXTMESSAGE:
            payload_size = mumble_proto__text_message__get_packed_size(data);
            break;
        case PACKET_USERREMOVE:
            payload_size = mumble_proto__user_remove__get_packed_size(data);
            break;
        case PACKET_USERSTATE:
            payload_size = mumble_proto__user_state__get_packed_size(data);
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
                mumble_proto__version__pack(data, packet_out.buffer + 6);
                break;
            case PACKET_UDPTUNNEL:
                memmove(packet_out.buffer + 6, data, length);
                break;
            case PACKET_AUTHENTICATE:
                mumble_proto__authenticate__pack(data, packet_out.buffer + 6);
                break;
            case PACKET_PING:
                mumble_proto__ping__pack(data, packet_out.buffer + 6);
                break;
            case PACKET_TEXTMESSAGE:
                mumble_proto__text_message__pack(data, packet_out.buffer + 6);
                break;
            case PACKET_USERREMOVE:
                mumble_proto__user_remove__pack(data, packet_out.buffer + 6);
                break;
            case PACKET_USERSTATE:
                mumble_proto__user_state__pack(data, packet_out.buffer + 6);
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
    lua_getfield(lua, -1, "threads");
    lua_pushnumber(lua, thread_id);
    lua_gettable(lua, -2);
    lua_getfield(lua, -1, "userthread");
    user_thread = (UserThread *)lua_touserdata(lua, -1);
    if (user_thread != NULL) {
        ev_io_stop(loop, w);
        free(user_thread);
    }
    lua_getfield(lua, -4, "Thread");
    lua_getfield(lua, -1, "_implFinish");
    lua_pushnumber(lua, thread_id);
    lua_call(lua, 1, 0);
}

void
script_stat_event(struct ev_loop *loop, ev_stat *w, int revents)
{
    ScriptStat *stat = (ScriptStat *)w;
    if (w->attr.st_ino == w->prev.st_ino && w->attr.st_mtime == w->prev.st_mtime) {
        return;
    }
    fprintf(stderr, "%s: reloaded %s\n", progname, stat->filename);
    lua_getglobal(lua, "piepan");
    lua_getfield(lua, -1, "_implLoadScript");
    lua_pushinteger(lua, stat->id);
    lua_call(lua, 1, 2);
    if (!lua_toboolean(lua, -2)) {
        fprintf(stderr, "%s: %s\n", progname, lua_tostring(lua, -1));
    }
    lua_settop(lua, 0);
}

static void
usage()
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
    fprintf(stderr, str, progname);
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
    ScriptStat *scripts = NULL;

    int socket_fd;
    struct sockaddr_in server_addr;

    ev_loop_main = EV_DEFAULT;
    SSLRead socket_watcher;
    ev_io user_thread_watcher;
    ev_timer ping_watcher;
    ev_signal signal_watcher;

    progname = argv[0];

    /*
     * Lua initialization
     */
    lua = luaL_newstate();
    if (lua == NULL) {
        fprintf(stderr, "%s: could not initialize Lua\n", progname);
        return 1;
    }
    luaL_openlibs(lua);
    api_init(lua);
    ret = 0;
    if (lua_load(lua, impl_reader, &ret, "piepan_impl", NULL) != 0) {
        fprintf(stderr, "%s: could not load piepan implementation\n", progname);
        return 1;
    }
    lua_call(lua, 0, 0);

    /*
     * Argument parsing
     */
    {
        int i;
        int show_help = 0;
        int show_version = 0;
        lua_getglobal(lua, "piepan");
        lua_getfield(lua, -1, "_implArgument");
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
                        progname, argv[i]);
                return 1;
            }
        }
        if (show_version) {
            printf("piepan %s (compiled on " __DATE__ " " __TIME__ ")\n",
                   PIEPAN_VERSION);
            return 0;
        }
        if (show_help) {
            usage();
            return 0;
        }
    }

    lua_settop(lua, 0);

    /*
     * Load user scripts
     */
    {
        int i;
        lua_settop(lua, 0);
        lua_getglobal(lua, "piepan");
        lua_getfield(lua, -1, "_implLoadScript");
        for (i = script_argc; i >= 0 && i < argc; i++) {
            lua_pushvalue(lua, -1);
            lua_pushstring(lua, argv[i]);
            lua_call(lua, 1, 2);
            if (lua_toboolean(lua, -2)) {
                if (developement_mode) {
                    ScriptStat *item = malloc(sizeof(ScriptStat));
                    if (item == NULL) {
                        fprintf(stderr, "%s: memory allocation error\n", progname);
                        return 1;
                    }
                    item->id = lua_tointeger(lua, -1);
                    item->filename = argv[i];
                    item->next = scripts;
                    scripts = item;
                    ev_stat_init(&item->ev, script_stat_event, item->filename, 0);
                    ev_stat_start(ev_loop_main, &item->ev);
                }
            } else {
                fprintf(stderr, "%s: %s\n", progname, lua_tostring(lua, -1));
            }
            lua_pop(lua, 2);
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
                    progname, opus_strerror(error));
            return 1;
        }
        opus_encoder_ctl(encoder, OPUS_SET_VBR(1));
        // TODO: set this to the server's max bitrate
        opus_encoder_ctl(encoder, OPUS_SET_BITRATE(50000));

        lua_settop(lua, 0);
    }

    /*
     * SSL initialization
     */
    SSL_library_init();

    ssl_context = SSL_CTX_new(SSLv23_client_method());
    if (ssl_context == NULL) {
        fprintf(stderr, "%s: could not create SSL context\n", progname);
        return 1;
    }

    if (certificate_file != NULL) {
        if (!SSL_CTX_use_certificate_chain_file(ssl_context, certificate_file) ||
                !SSL_CTX_use_PrivateKey_file(ssl_context, key_file,
                                                SSL_FILETYPE_PEM) ||
                !SSL_CTX_check_private_key(ssl_context)) {
            fprintf(stderr, "%s: could not load certificate and/or key file\n",
                    progname);
            return 1;
        }
    }

    /*
     * Socket initialization and connection
     */
    socket_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (socket_fd < 0) {
        fprintf(stderr, "%s: could not create socket\n", progname);
        return 1;
    }

    memset(&server_addr, 0, sizeof(server_addr));
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(port);

    server_host = gethostbyname(server_host_str);
    if (server_host == NULL || server_host->h_addr_list[0] == NULL ||
            server_host->h_addrtype != AF_INET) {
        fprintf(stderr, "%s: could not parse server address\n", argv[0]);
        return 1;
    }
    memmove(&server_addr.sin_addr, server_host->h_addr_list[0],
            server_host->h_length);

    ret = connect(socket_fd, (struct sockaddr *) &server_addr,
                  sizeof(server_addr));
    if (ret != 0) {
        fprintf(stderr, "%s: could not connect to server\n", progname);
        return 1;
    }

    ssl = SSL_new(ssl_context);
    if (ssl == NULL) {
        fprintf(stderr, "%s: could not create SSL object\n", progname);
        return 1;
    }

    if (SSL_set_fd(ssl, socket_fd) == 0) {
        fprintf(stderr, "%s: could not set SSL file descriptor\n", progname);
        return 1;
    }

    if (SSL_connect(ssl) != 1) {
        fprintf(stderr, "%s: could not create secure connection\n", progname);
        return 1;
    }

    /*
     * User thread pipe
     */
    if (pipe(user_thread_pipe) != 0) {
        fprintf(stderr, "%s: could not create user thread pipe\n", progname);
        return 1;
    }

    /*
     * Send initial packets
     */
    {
        MumbleProto__Version version = MUMBLE_PROTO__VERSION__INIT;
        MumbleProto__Authenticate auth = MUMBLE_PROTO__AUTHENTICATE__INIT;
        FILE *file;
        char buffer[TOKEN_BUFFER_SIZE];
        struct {
            int count;
            char *arr[MAX_TOKENS];
        } tokens = {0};

        auth.has_opus = true;
        auth.opus = true;
        auth.username = username;
        if (password_file != NULL) {
            file = fopen(password_file, "r");
            if (file == NULL) {
                fprintf(stderr, "%s: could open password file for reading\n",
                        progname);
                return 1;
            }
            if (fgets(buffer, sizeof(buffer), file) == NULL) {
                fprintf(stderr, "%s: could not read password from file\n",
                        progname);
                fclose(file);
                return 1;
            }
            rnltrim(buffer, strlen(buffer));
            auth.password = buffer;
            fclose(file);
        }
        if (token_file != NULL) {
            file = fopen(token_file, "r");
            if (file == NULL) {
                fprintf(stderr, "%s: could open token file for reading\n",
                        progname);
                return 1;
            }
            while (tokens.count < MAX_TOKENS) {
                if (fgets(buffer, sizeof(buffer), file) == NULL) {
                    break;
                }
                if (buffer[0] == '\n') {
                    continue;
                }
                rnltrim(buffer, strlen(buffer));
                tokens.arr[tokens.count++] = strdup(buffer);
                tokens.count++;
            }
            if (tokens.count > 0) {
                auth.n_tokens = tokens.count;
                auth.tokens = tokens.arr;
            }
            fclose(file);
        }

        version.has_version = true;
        version.version = 1 << 16 | 2 << 8 | 4; // 1.2.4
        version.release = "Unknown";
        version.os = "piepan";
        version.os_version = PIEPAN_VERSION;

        sendPacket(PACKET_VERSION, &version);
        sendPacket(PACKET_AUTHENTICATE, &auth);

        while (tokens.count > 0) {
            free(tokens.arr[--tokens.count]);
        }
    }

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
    lua_getfield(lua, -1, "_implOnDisconnect");
    if (lua_isfunction(lua, -1)) {
        lua_newtable(lua);
        lua_call(lua, 1, 0);
    }

    SSL_shutdown(ssl); // TODO:  sigpipe is triggered here if connection breaks
    close(socket_fd);
    lua_close(lua);

    if (developement_mode) {
        ScriptStat *item = scripts;
        while (item != NULL) {
            ScriptStat *next = item->next;
            free(item);
            item = next;
        }
    }

    return 0;
}
