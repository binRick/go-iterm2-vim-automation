#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

jo now=$(date +%s) abc=123 | ./gen_seq.sh
