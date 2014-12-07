#!/bin/bash
set -e
echo "$1"
wget -q -O - "http://gdata.youtube.com/feeds/api/videos/$1?v=2&alt=jsonc" | \
  jshon -Q -e data -e title -u -p -e duration -u -p -e thumbnail -e hqDefault -u
