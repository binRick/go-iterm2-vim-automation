#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )


cmd='./new_tab.sh|jq'
exec nodemon -w . -e go,sh -x sh -- -c "$cmd||true"
