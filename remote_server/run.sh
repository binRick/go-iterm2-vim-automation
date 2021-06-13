#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ARGS="$@"


cmd="command go run . $ARGS"
exec $cmd
