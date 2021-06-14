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

ACTION=list
#DIR="%2ftmp"
#DIR="%2FUsers%2Frick%2FDesktop%2Fgo-iterm2-vim-automation%2Flocal_server%0A"
#CMD="pwd"
#cmd="curl $CURL_ARGS http://$CURL_HOST:$CURL_PORT/api/iterm2/new_tab?hostname=localhost\&directory=${DIR}&cmd=$CMD"

#>&2 echo -e "$cmd"
eval $cmd
