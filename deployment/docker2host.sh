#!/bin/bash
set -e

PORTS=($(echo "$1" | grep -oP '^\d+(:\d+)?$' | sed -e 's/:/ /g'))
if [ -z $PORTS ]; then
    cat <<EOF
Usage:
$(basename "$0") SRC[:DEST]
    SRC: will be the port accessible inside the container
    DEST:
        the connection will be redirected to this port on the host.
        if not specified, the same port as SRC will be used
EOF
    exit 1
fi

SOURCE=${PORTS[0]}
DEST=${PORTS[1]-${PORTS[0]}}

SOCKFILE="$XDG_RUNTIME_DIR/forward-docker2host-${SOURCE}_$DEST.sock"
# socat UNIX-LISTEN:"$SOCKFILE",fork TCP:127.0.0.1:$DEST &  ## For TCP
socat UNIX-LISTEN:"$SOCKFILE",fork TCP6:[::]:$DEST &
nsenter -U --preserve-credentials -n -t $(cat "$XDG_RUNTIME_DIR/docker.pid") -- socat TCP4-LISTEN:$SOURCE,reuseaddr,fork "$SOCKFILE" &
echo forwarding $SOURCE:$DEST... use ctrl+c to quit

sleep 365d
