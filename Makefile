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
	rm -f readme.html

.PHONY: clean
