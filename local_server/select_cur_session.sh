#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
INITIAL_QUERY="$@"
source bin/ansi
export PATH="$PATH:$(pwd)/bin"

RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case --max-depth $MAX_DEPTH"
cmd="$(command -v fzf) --bind 'change:reload:$RG_PREFIX {q} || true' --ansi --disabled --query '$INITIAL_QUERY' --header 'Source Code Search' --header-lines=0 --height=80% --layout=reverse"

FZF_TITLE="\
$(ansi --white --bg-black --bold --underline "ctrl+m:INFO")"

tab_new_bat_cmd() {
  cmd="./new_bat_tab.sh {}"
  echo -e "$cmd"
}
tab_new_vim_cmd() {
  cmd="./new_vim_tab.sh {}"
  echo -e "$cmd"
}
tab_new_cmd() {
  cmd="./new_ls_tab.sh"
  echo -e "$cmd"
}



tab_info_cmd() {
  cmd1="ansi --bold --yellow   \"v:   Open File with Vim editor in new tab to the right\""
  cmd2="ansi --italic --magenta  \"r:  Run file as executable in new tab to the right\""
  cmd3="ansi --italic --cyan  \"n:  Run nodemon on executable in new tab to the right\""
  cmd4="ansi --italic --cyan  \"b:  Run bat on file in new tab to the right\""
  echo -e "$cmd1;$cmd2;$cmd3;$cmd4;"
}

fzf_cmd="fzf --height 100% --layout reverse --info inline --border \
    --ansi \
    --header='$FZF_TITLE' \
    --preview './test_session.sh {}' \
    --border=sharp \
    -m \
    --bind 'ctrl-t:preview:$(tab_new_cmd)'\
    --bind 'ctrl-m:preview:$(tab_info_cmd)'\
    --bind 'ctrl-v:preview:$(tab_new_vim_cmd)'\
    --bind 'ctrl-b:preview:$(tab_new_bat_cmd)'\
    --preview-window right,80%,border-vertical:sharp \
    --color 'fg:#bbccdd,fg+:#ddeeff,bg:#334455,preview-bg:#223344,border:#778899'"

cmd="./get_cur_sessions.sh|$fzf_cmd"

>&2 echo -e "$cmd"
eval $cmd
