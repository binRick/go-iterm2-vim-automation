#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
./get_iterm_sessions.sh|jq '.session' -Mrc
#./list_vims.sh|jq '.[].Session' -Mrc
