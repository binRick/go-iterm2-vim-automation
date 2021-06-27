#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
SESSION="$1"


REMOTE_PORT=48923

CURL_ARGS="-s"
CURL_HOST=127.0.0.1
CURL_PORT=$REMOTE_PORT

cmd="curl $CURL_ARGS http://$CURL_HOST:$CURL_PORT/api/iterm2/dump/session/$SESSION"

>&2 echo -e "$cmd"
eval $cmd | base64 -d
