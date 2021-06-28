#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )


nodemon -I -w . -e go -x sh -- -c "./run.sh||true"
