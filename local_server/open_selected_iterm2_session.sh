#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
SESSION="$1"
TAB="$2"


REMOTE_PORT=48923

CURL_ARGS="-s"
CURL_HOST=127.0.0.1
CURL_PORT=$REMOTE_PORT

cmd="curl $CURL_ARGS http://$CURL_HOST:$CURL_PORT/api/iterm2/open/session/$SESSION/tab/$TAB"

>&2 echo -e "$cmd"
eval $cmd
