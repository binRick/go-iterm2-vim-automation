#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

ps ax|grep -v 'grep ' | grep 15229:127.0.0.1:|cut -d ' ' -f1|sort -u
