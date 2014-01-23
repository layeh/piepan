CFLAGS = `pkg-config --libs --cflags libssl lua libprotobuf-c opus vorbis vorbisfile` -lev -pthread

LUAFILES = src/impl/piepan.lua \
           src/impl/meta.lua \
           src/impl/scripts.lua \
           src/impl/timer.lua \
           src/impl/thread.lua \
           src/impl/user.lua \
           src/impl/channel.lua \
           src/impl/events.lua \
           src/impl/functions.lua

piepan: src/piepan.c src/piepan.h src/util.c src/handlers.c src/api.c \
	proto/Mumble.o src/piepan_impl.c
	$(CC) $(CFLAGS) -o $@ src/piepan.c proto/Mumble.o

proto/Mumble.o: proto/Mumble.proto
	protoc-c --c_out=. proto/Mumble.proto
	$(CC) -c -I. -o proto/Mumble.o proto/Mumble.pb-c.c

src/piepan_impl.c: src/piepan_impl.luac
	xxd -i src/piepan_impl.luac src/piepan_impl.c

src/piepan_impl.luac: src/piepan_impl.lua
	luac -o src/piepan_impl.luac src/piepan_impl.lua

src/piepan_impl.lua: $(LUAFILES)
	cat $(LUAFILES) > $@

readme.html: README.md
	echo '<!DOCTYPE html>' > readme.html
	echo '<html>' >> readme.html
	echo '<head>' >> readme.html
	echo '<meta charset="utf-8" />' >> readme.html
	echo '<title>piepan: a bot framework for Mumble</title>' >> readme.html
	echo '<style type="text/css">' >> readme.html
	echo 'body {font-family:sans-serif;max-width:900px;margin:0 auto;padding:10px}' >> readme.html
	echo '</style>' >> readme.html
	echo '</head>' >> readme.html
	echo '<body>' >> readme.html
	markdown README.md >> readme.html
	echo '</body>' >> readme.html
	echo '</html>' >> readme.html

clean:
	rm -f piepan
	rm -f proto/Mumble.o proto/Mumble.pb-c.c proto/Mumble.pb-c.h
	rm -f src/piepan_impl.c src/piepan_impl.luac src/piepan_impl.lua
	rm -f readme.html

.PHONY: clean
