#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
VIM_FILE="$@"

CURL_ARGS="-s"
CURL_HOST=127.0.0.1
CURL_PORT=15223


cmd="curl $CURL_ARGS http://$CURL_HOST:$CURL_PORT/list"

eval $cmd
