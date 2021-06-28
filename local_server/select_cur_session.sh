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

open_selected_iterm2_session() {
  cmd="./open_selected_iterm2_session.sh {}"
  echo -e "$cmd"
}

tab_info_cmd() {
  cmd1="ansi --bold --yellow     \"ctrl+o:  Switch to selected tab\""
  echo -e "$cmd1;"
}

items_cmd="./get_cur_sessions.sh"

#    --bind='ctrl-r:reload($items_cmd)' \
fzf_cmd="fzf --height 100% --layout=reverse --info inline --border \
    --ansi \
    --header='$FZF_TITLE' \
    --preview='passh ./get_session_uuid_info.sh {}' \
    --header-lines=0 \
    --border=sharp \
    -m \
    --bind 'ctrl-h:preview:$(tab_info_cmd)'\
    --bind 'ctrl-o:preview:$(open_selected_iterm2_session)'\
    --preview-window=right,80%,border-vertical:sharp \
    --color 'fg:#bbccdd,fg+:#ddeeff,bg:#334455,preview-bg:#223344,border:#778899'"

cmd="$items_cmd | $fzf_cmd"

#>&2 echo -e "$cmd"
selected_session_id="$(eval $cmd)"
cmd="./open_selected_iterm2_session.sh $selected_session_id"
eval $cmd
