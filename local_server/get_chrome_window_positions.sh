#!/bin/bash
set -e
cd $( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

chrome-cli list windows|cut -d']' -f1|cut -d'[' -f2|sort -u|while read -r w; do 
  echo -ne window=$w,; chrome-cli position -w $w
done | sed 's/[[:space:]]//g'|tr ',' ' '|tr ':' '=' | while read p; do
  jo $p
done
