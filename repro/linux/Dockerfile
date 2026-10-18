FROM ubuntu:24.04
COPY --from=golang:1.27-bookworm /usr/local/go /usr/local/go
ENV PATH=/usr/local/go/bin:/root/go/bin:$PATH GOTOOLCHAIN=local
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates \
    fonts-dejavu fonts-inter fonts-noto-cjk fontconfig libgtk-3-0t64 libegl1 libegl-mesa0 libgl1-mesa-dri libepoxy0 libxtst6 \
    xvfb xdotool imagemagick x11-utils dbus-x11 && rm -rf /var/lib/apt/lists/*

