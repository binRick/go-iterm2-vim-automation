#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

set +e
while read -r pid; do
  kill $pid
done < <(./get_ssh_pids.sh)
