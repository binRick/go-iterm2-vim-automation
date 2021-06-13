#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ARGS="$@"
BDIR="$(cd ../../../. && pwd)"


cmd="nodemon -w . \
--signal SIGINT \
--delay .2 \
-w $BDIR/utils \
-w $BDIR/blockers \
-w $BDIR/blocker \
 -e go -x './run.sh $ARGS||true'"


#echo -e "$cmd"
eval $cmd
