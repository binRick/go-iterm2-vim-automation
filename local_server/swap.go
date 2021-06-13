package main

var EXTRACT_ENV_JSONS_CMD = `cat /proc/*/environ 2>/dev/null | tr '\0' '\n' | grep ^ITERM_PROFILE_VARS_ENCODED=|cut -d'=' -f2-100|while read -r h; do echo -e $h|base64 -d 2>/dev/null|grep '^{'|jq -Mrc;done`

type VimCrash struct {
	SwapFileName string
	VimFilePath  string
	Modified     bool
	User         string
	Hostname     string
	Pid          int
	Running      bool

	DateString string

	Length int
}
