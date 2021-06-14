#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
VIM_FILE="$1"

REMOTE_PORT=48923
LOCAL_PORT=19344
REMOTE_LISTEN_HOST=127.0.0.1

CURL_ARGS="-s"
CURL_HOST=127.0.0.1
CURL_PORT=$REMOTE_PORT

cmd="curl $CURL_ARGS http://$CURL_HOST:$CURL_PORT/api/iterm2/info"

#>&2 echo -e "$cmd"
eval $cmd
