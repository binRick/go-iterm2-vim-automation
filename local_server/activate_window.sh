#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ARGS="$@"

REMOTE_PORT=48923
LOCAL_PORT=19344
REMOTE_LISTEN_HOST=127.0.0.1

CURL_ARGS="-s"
CURL_HOST=127.0.0.1
CURL_PORT=$REMOTE_PORT

ACTION=test
WINDOW=6
SESSION=DE3015A1-67F0-4A72-934A-86B530A2F0EE
TAB=0

cmd="curl $CURL_ARGS http://$CURL_HOST:$CURL_PORT/api/iterm2/$ACTION/window/$WINDOW/session/$SESSION/tab/$TAB"

>&2 echo -e "$cmd"
eval $cmd
