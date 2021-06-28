#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

ARGS="$@"
cmd="./run.sh $ARGS"

cmd="command nodemon --signal SIGKILL -w . -e sh,go -x sh -c -- '$cmd||true'"

exec $cmd
