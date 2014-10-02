LUAFILES = impl/piepan.lua \
           impl/internal.lua \
           impl/scripts.lua \
           impl/timer.lua \
           impl/thread.lua \
           impl/user.lua \
           impl/channel.lua \
           impl/events.lua \
           impl/permissions.lua \
           impl/audio.lua \
           impl/functions.lua

impl/piepan_impl.go: impl/piepan_impl.lua
	echo 'package impl' > $@
	echo 'const Piepan = `' >> $@
	cat $< >> $@
	echo '`' >> $@

impl/piepan_impl.lua: $(LUAFILES)
	cat $(LUAFILES) > $@

readme.html: README.md
	echo '<!DOCTYPE html>' > $@
	echo '<html>' >> $@
	echo '<head>' >> $@
	echo '<meta charset="utf-8" />' >> $@
	echo '<title>piepan: a bot framework for Mumble</title>' >> $@
	echo '<style type="text/css">' >> $@
	echo 'body {font-family:sans-serif;margin:0 auto;padding: 0 10px}' >> $@
	echo '</style>' >> $@
	echo '</head>' >> $@
	echo '<body>' >> $@
	markdown README.md >> $@
	echo '</body>' >> $@
	echo '</html>' >> $@

clean:
	rm -f piepan
	rm -f impl/piepan_impl.lua impl/piepan_impl.go
	rm -f readme.html

.PHONY: clean
