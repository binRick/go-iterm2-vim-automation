#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )


while IFS=: read -r pid iterm_id; do
  if ! \ps -p $pid >/dev/null 2>&1; then 
    continue
  else
    true
  fi
    wwttpp="$(echo -e "$iterm_id"|cut -d: -f1|tr '[a-z]' ':')"
    w="$(echo -e "$wwttpp"|cut -d: -f2)"
    t="$(echo -e "$wwttpp"|cut -d: -f3)"
    p="$(echo -e "$wwttpp"|cut -d: -f4)"
    iterm_uuid="$(echo -e "$iterm_id"|cut -d: -f2)"
    [[ "$iterm_uuid" == "" ]] && continue
    [[ "$w" == "" ]] && continue
    cmd="jo window=$w tab=$t p=$p session=$iterm_uuid pid=$pid"
    eval $cmd
done < <(cat /tmp/iterm_session.log)
