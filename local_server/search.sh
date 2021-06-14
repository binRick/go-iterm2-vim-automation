#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
INITIAL_QUERY="$@"
source bin/ansi
export PATH="$PATH:$(pwd)/bin"

RG_PREFIX="rg --column --line-number --no-heading --color=always --smart-case --max-depth $MAX_DEPTH"
cmd="$(command -v fzf) --bind 'change:reload:$RG_PREFIX {q} || true' --ansi --disabled --query '$INITIAL_QUERY' --header 'Source Code Search' --header-lines=0 --height=80% --layout=reverse"

FZF_TITLE="\
$(ansi --white --bg-black --bold --underline "control+:   ")\
$(ansi --green --italic "t-Tab"),\
"

tab_new_vim_cmd() {
  cmd="./new_vim_tab.sh {}"
  echo -e "$cmd"
}
tab_new_cmd() {
  cmd="./new_ls_tab.sh"
  echo -e "$cmd"
}

tab_info_cmd() {
  cmd1="ansi --bold --yellow   \"Session:   yyyyyyyyyyyyy\""
  cmd2="ansi --italic --magenta  \"Window:    yyyyyyyyyyyyy\""
  cmd3="ansi --italic --cyan     \"Tab:       yyyyyyyyyyyyy\""
  cmd4="ansi --italic --white    \"Pid:       213\""
  echo -e "$cmd1;$cmd2;$cmd3;$cmd4;"
}

cmd="fzf --height 100% --layout reverse --info inline --border \
    --bind 'change:reload:$RG_PREFIX {q} || true' \
    --ansi --query '$INITIAL_QUERY' --disabled \
    --header='$FZF_TITLE' \
    --preview 'bat --color=always {}' \
    --border=sharp \
    -m \
    --bind 'ctrl-t:preview:$(tab_new_cmd)'\
    --bind 'ctrl-i:preview:$(tab_info_cmd)'\
    --bind 'ctrl-v:preview:$(tab_new_vim_cmd)'\
    --preview-window right,80%,border-vertical:sharp \
    --color 'fg:#bbccdd,fg+:#ddeeff,bg:#334455,preview-bg:#223344,border:#778899'"

>&2 echo -e "$cmd"
eval $cmd
