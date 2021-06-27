#!/usr/bin/env bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
INITIAL_QUERY="$@"
source bin/ansi
export PATH="$PATH:$(pwd)/bin"

RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case --max-depth $MAX_DEPTH"
cmd="$(command -v fzf) --bind 'change:reload:$RG_PREFIX {q} || true' --ansi --disabled --query '$INITIAL_QUERY' --header 'Source Code Search' --header-lines=0 --height=80% --layout=reverse"

FZF_TITLE="\
$(ansi --white --bg-black --bold --underline "ctrl+h:HELP")"

info_cmd() {
  cmd="echo ./new_vim_tab.sh {}"
  echo -e "$cmd"
}



help_cmd() {
  cmd1="ansi --italic --cyan  \"i:  VIM Session Info\""
  echo -e "$cmd1;"
}

items_cmd="./get_vim_sessions.sh"

fzf_cmd="fzf --height 100% --layout=reverse --info inline --border \
    --ansi \
    --header='$FZF_TITLE' \
    --preview='echo passh ./test_session.sh {}' \
    --header-lines=0 \
    --border=sharp \
    -m \
    --bind 'ctrl-h:preview:$(help_cmd)'\
    --bind 'ctrl-i:preview:$(info_cmd)'\
    --preview-window=right,80%,border-vertical:sharp \
    --color 'fg:#bbccdd,fg+:#ddeeff,bg:#334455,preview-bg:#223344,border:#778899'"

cmd="$items_cmd | $fzf_cmd"

eval $cmd
