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
	echo '<!DOCTYPE html>' > readme.html
	echo '<html>' >> readme.html
	echo '<head>' >> readme.html
	echo '<meta charset="utf-8" />' >> readme.html
	echo '<title>piepan: a bot framework for Mumble</title>' >> readme.html
	echo '<style type="text/css">' >> readme.html
	echo 'body {font-family:sans-serif;margin:0 auto;padding: 0 10px}' >> readme.html
	echo '</style>' >> readme.html
	echo '</head>' >> readme.html
	echo '<body>' >> readme.html
	markdown README.md >> readme.html
	echo '</body>' >> readme.html
	echo '</html>' >> readme.html

clean:
	rm -f piepan
	rm -f impl/piepan_impl.lua impl/piepan_impl.go
	rm -f readme.html

.PHONY: clean
