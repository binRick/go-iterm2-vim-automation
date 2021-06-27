#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

./list_vims.sh | jq '.[].Session' -Mrc | tac
