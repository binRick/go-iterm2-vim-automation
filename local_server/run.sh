#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
ARGS="$@"
REMOTE_PORT=48923
LOCAL_PORT=19344
REMOTE_LISTEN_HOST=127.0.0.1
rem_cmd="netstat -alntp|grep LISTEN|grep sshd|grep $REMOTE_LISTEN_HOST:$REMOTE_PORT"
[[ -f ~/.ssh/config.none ]] || touch ~/.ssh/config.none
REMOTE_SERVER=$__REMOTE_SERVER
REMOTE_USER=$__REMOTE_USER

export VIM_LOCAL_PORT=15229
VIM_REMOTE_HOST=127.0.0.1
VIM_REMOTE_PORT=15223

PORT_FORWARDS="-R $REMOTE_PORT:$REMOTE_LISTEN_HOST:$LOCAL_PORT -L $VIM_LOCAL_PORT:$VIM_REMOTE_HOST:$VIM_REMOTE_PORT"

base_ssh_cmd="command ssh -oStrictHostKeyChecking=no -oLogLevel=ERROR -F ~/.ssh/config.none -oControlMaster=no $PORT_FORWARDS $REMOTE_USER@$REMOTE_SERVER"
ssh_cmd="$base_ssh_cmd $rem_cmd"
sleep_ssh_cmd="$base_ssh_cmd -N"

ENV_FILE=~/.iterm_profile.json


props="local_port=$LOCAL_PORT remote_port=$REMOTE_PORT remote_listen_host=$REMOTE_LISTEN_HOST remote_user=$REMOTE_USER remote_server=$REMOTE_SERVER ts=$(date +%s) pid=$$ hostname=$(hostname -f) user=$(id -u --name)"
env_file_cmd="echo -e \"$props\""
env_file_cmd="jo $props"

eval $env_file_cmd > $ENV_FILE


MAX_RETRIES=10
CUR_RETRIES=0
SLEEP_RETRY=.2
>&2 echo -e "$ssh_cmd"
set +e
while :; do
  if ! eval timeout 3 $ssh_cmd; then
    if [[ "$CUR_RETRIES" -ge "$MAX_RETRIES" ]]; then
      echo SSH Port Forward Failed
      exit 1
    fi
  else
    break
  fi
  CUR_RETRIES="$(($CUR_RETRIES+1))"
  sleep $SLEEP_RETRY
done

set -e

ssh_pids=""
kill_ssh(){
  echo Killing ssh!

  echo -e "$ssh_pids"
}

eval $sleep_ssh_cmd &

ssh_pids="$ssh_pids $!"

trap kill_ssh EXIT

#bin=.sb
#go build -o $bin .
#cmd="./$bin $@"
cmd="go run . $@"
#killall $bin 2>/dev/null ||true
eval $cmd
