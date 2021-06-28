package main

import (
	"github.com/mgutz/ansi"
	"mrz.io/itermctl"
)

var (
	MONITOR_CONTROL_SEQUENCE = true
	CONTROL_SEQUENCE_REGEX   = `vim-iterm`
	CONTROL_SEQUENCE_NAME    = `vim-iterm`
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
	_conn     *itermctl.Connection
	_app      *itermctl.App
	ctrl_seq1 = itermctl.NewCustomControlSequenceEscaper(CONTROL_SEQUENCE_NAME)
	err_msg   = ansi.ColorFunc("red")
)
