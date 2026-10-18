#!/bin/bash
# shot.sh <before|after> <zoom>: the demo in an X11 window, captured.
v=$1; z=$2
Xvfb :9 -screen 0 1200x800x24 >/dev/null 2>&1 & sleep 1
export DISPLAY=:9 HOME=/tmp
eval $(dbus-launch --sh-syntax)
ZOOM=$z /bins/zoomdemo-$v >/tmp/app.log 2>&1 & sleep 5
W=$(xdotool search --onlyvisible --name "Zoom demo" | head -1)
xdotool windowsize $W 1200 800 windowmove $W 0 0; sleep 2
import -window "$W" /out/linux-$v-$z.png
