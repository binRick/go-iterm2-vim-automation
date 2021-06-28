#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

vimfile=/etc/passwd

jo type=vim file=$vimfile | ./gen_seq.sh
