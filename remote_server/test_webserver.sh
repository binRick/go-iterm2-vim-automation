#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ARGS="$@"

CURL_ARGS="-v"
CURL_HOST=127.0.0.1
CURL_PORT=15223


cmd="curl $CURL_ARGS http://$CURL_HOST:$CURL_PORT/"

>&2 echo -e "$cmd"
eval $cmd
