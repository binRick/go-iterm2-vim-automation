#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
session_id="$1"


j="$(./get_iterm_sessions.sh|grep \"$session_id\"|jq  -Mrc)"

window_id="$(echo -e "$j"|jq '.window')"
tab_id="$(echo -e "$j"|jq '.tab')"

REMOTE_PORT=48923
LOCAL_PORT=19344
REMOTE_LISTEN_HOST=127.0.0.1

CURL_ARGS="-s"
CURL_HOST=127.0.0.1
CURL_PORT=$REMOTE_PORT

cmd="curl $CURL_ARGS http://$CURL_HOST:$CURL_PORT/api/iterm2/test/window/${window_id}/session/${session_id}/tab/${tab_id}"

#>&2 echo -e "$cmd"
eval $cmd|base64 -d
