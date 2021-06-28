#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

export NOTIFY_INTERVAL=${NOTIFY_INTERVAL:-5}

CONTROL_SEQUENCE_NAME=${CONTROL_SEQUENCE_NAME:-vim-iterm}
CONTROL_SEQUENCE_PREFIX="${1:-test-seq}"
CONTROL_SEQUENCE_DATA="${2:--}"
DAT="blah123"
if [[ "$CONTROL_SEQUENCE_DATA" == "-" ]]; then
  DAT=
  while read -r line; do
    DAT="$DAT\n$line"
  done
fi

cwd="$(pwd)"

jo_cmd="jo dat='$(echo -e "$DAT"|jq -Mrc)' ts=$(date +%s) hostname='$(command hostname -f)' user=$(id -nu) cwd='$cwd' ITERM_SESSION_ID='$ITERM_SESSION_ID' ITERM_PROFILE='$ITERM_PROFILE' pid=$$ SSH_CONNECTION='$SSH_CONNECTION'"
encode_cmd="base64 -w0"
DAT="$(eval $jo_cmd | jq -Mrc | $encode_cmd)"

msg="\033]1337;Custom=id=${CONTROL_SEQUENCE_NAME}:${CONTROL_SEQUENCE_PREFIX}:${DAT}\a"
echo -e "$msg"

