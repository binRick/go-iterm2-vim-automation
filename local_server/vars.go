package main

import (
	"os"
	"strconv"
)

var (
	VIM_QUIT_TEXT = "q"
)
var (
	KEYSTROKE_ACTION_DEBUG_DETAILED = false
	KEYSTROKE_ACTION_DEBUG_RESULT   = true
	KEYSTROKE_ACTION_SHOW_ALERT     = false
)
var (
	KEYSTROKE_ACTION_SELECT_NEW_SESSION   = true
	KEYSTROKE_ACTION_SELECT_NEW_TAB       = false
	KEYSTROKE_ACTION_QUIT_OLD_SESSION_VIM = true
)

var (
	DEBUG_RESULT     = true
	VIM_LOCAL_PORT   uint
	DEBUG_KEYSTROKES = false
)

var MATCHED_VIM_SWAP_STRINGS = []string{
	`Found a swap file by the name`,
}

func init() {
	_port, err := strconv.ParseInt(os.Getenv(`VIM_LOCAL_PORT`), 10, 0)
	if err == nil {
		VIM_LOCAL_PORT = uint(_port)
	}
}
