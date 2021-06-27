#!/usr/bin/env bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
source ansi
session_id="$1"



list_offspring(){
  tp=`pgrep -P $1`          #get childs pids of parent pid
  for i in $tp; do          #loop through childs
    if [ -z $i ]; then      #check if empty list
      exit                  #if empty: exit
    else                    #else
      echo -n "$i "         #print childs pid
      if [[ "$i" != "" ]]; then
        list_offspring $i     #call list_offspring again with child pid as the parent
      fi
    fi;
  done
}





grep ":$session_id$" /tmp/iterm_session.log -q

#TERM_HEMGHT="ansi --report-window-chars|cut -d, -f2|grep '^[0-9]'"

pid="$(grep ":$session_id$" /tmp/iterm_session.log |cut -d: -f1)"
pids="$(list_offspring $pid|tr '\n' ' '|tr ' ' '\n'|grep '^[0-9]'|tr '\n' ' ')"
info="$(grep ":$session_id$" /tmp/iterm_session.log |cut -d: -f2)"


wwttpp="$(echo -e "$info"|cut -d: -f1|tr '[a-z]' ':')"
w="$(echo -e "$wwttpp"|cut -d: -f2)"
t="$(echo -e "$wwttpp"|cut -d: -f3)"
p="$(echo -e "$wwttpp"|cut -d: -f4)"


#echo -e "$w=>$t=>$p"

j="$(./get_iterm_sessions.sh|grep \"$session_id\"|jq  -Mrc)"

window_id="$(echo -e "$j"|jq '.window')"
tab_id="$(echo -e "$j"|jq '.tab')"

REMOTE_PORT=48923
LOCAL_PORT=19344
REMOTE_LISTEN_HOST=127.0.0.1

CURL_ARGS="-s"
CURL_HOST=127.0.0.1
CURL_PORT=$REMOTE_PORT

cmd="curl $CURL_ARGS http://$CURL_HOST:$CURL_PORT/api/iterm2/test/window/${window_id}/session/${session_id}/tab/${tab_id}"

#>&2 echo -e "$cmd"
stdout="$(eval $cmd|base64 -d)"

#echo -e "$stdout"
#exit

render_stdout(){
  msg="$(ansi --yellow --italic "$stdout")"
  echo -e "$msg"
}

qty_lines="$(./get_session_lines_qty.sh $session_id)"

pids_report="$(ps -p $(echo -e "$pids"|tr ' ' ',') -O comm 2>/dev/null||true)"
pids_qty="$(($(echo -e "$pids_report"|wc -l)-1))"

ansi -n --green --bold "| Session $(ansi --yellow --italic "$session_id")"
ansi -n --green --bold " > Window $(ansi --yellow --italic "#$w")"
ansi -n --green --bold " > Tab $(ansi --yellow --italic "#$t") "
#ansi --reset
echo -ne "\n"
ansi -n --green --bold "| PID $(ansi --yellow --italic "$pid")"

if [[ "$pids_qty" -gt 0 ]]; then
  ansi -n --green --bold " > $(ansi --white --bg-black --bold $pids_qty) PIDs $(ansi --yellow --italic "$pids")"
fi

echo -ne "\n"
ansi -n --green --bold "| Session Lines Qty: $(ansi --yellow --italic "$qty_lines")"

if [[ "$pids_qty" -gt 0 ]]; then
  echo -e "\n\n$(ansi --white --bg-black --bold $pids_qty) Pids Activity:"
  ansi -n --cyan --bg-black --italic "$pids_report"
fi
echo -e "\n"

render_stdout
exit
