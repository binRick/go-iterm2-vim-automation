package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"dev.local/types"
	"github.com/k0kubun/pp"
	"github.com/pterm/pterm"
	"mrz.io/itermctl"
)

var test_msg = `  E325: ATTENTION
Found a swap file by the name ".xxx.go.swp"
owned by: root   dated: Sun Jun 13 10:55:33 2021
file name: /home/xxxxxxxxxxxxxxx/xxx.go
modified: no
user name: root   host name: host.hostname.domain
process ID: 224351 (STILL RUNNING)
While opening file "xxx.go"`

func init() {
	test_parse_crash_msg()
}

func get_iterm_profile_from_vim_crash(crash *VimCrash) (*types.ItermProfile, error) {
	p := types.ItermProfile{}
	if len(crash.VimFilePath) < 1 {
		err_msg := fmt.Sprintf("Vim crash has an invalid vim file path of '%s'", crash.VimFilePath)
		return &p, errors.New(err_msg)
	}

	vims, err := get_remote_vms(VIM_LOCAL_PORT)
	F(err)
	for _, vim := range *vims {
		if false {
			pp.Println(vim)
		}
		if vim.VimFilePath == crash.VimFilePath {
			return &vim, nil
		}
	}

	if len(p.Session) < 1 {
		err_msg := fmt.Sprintf("Failed to acquire iterm profile from vim crash for file '%s'", crash.VimFilePath)
		F(errors.New(err_msg))
	}
	return &p, nil
}

func test_parse_crash_msg() {
	//_ = parse_crash_msg(test_msg)
}

func parse_crash_msg(msg string) *VimCrash {
	c := VimCrash{
		Length:   len(msg),
		Modified: false,
		Running:  false,
	}
	for index, line := range strings.Split(msg, "\n") {
		if strings.Contains(line, `dated: `) {
			c.DateString = strings.Split(line, `dated: `)[1]
		}
		if strings.Contains(line, `host name: `) {
			c.Hostname = strings.Split(line, `host name: `)[1]
		}
		if strings.HasPrefix(line, `process ID: `) {
			_pid := strings.Split(line, ` `)[2]
			__pid, err := strconv.ParseInt(_pid, 10, 0)
			F(err)
			c.Pid = int(__pid)
			if strings.Contains(line, `STILL RUNNING`) {
				c.Running = true
			}
		}
		if strings.HasPrefix(line, `user name: `) {
			c.User = strings.Split(line, ` `)[2]
		}
		if strings.HasPrefix(line, `modified: `) {
			Modified := strings.Split(line, ` `)[1]
			if Modified == `yes` {
				c.Modified = true
			}
		}
		if strings.Contains(line, `file name: `) {
			c.VimFilePath = strings.Split(line, `file name: `)[1]
		}
		if strings.HasPrefix(line, `Found a swap file by the name`) {
			c.SwapFileName = strings.Split(strings.Split(line, ` `)[7], `"`)[1]
		}
		if false {
			fmt.Println(
				pp.Sprintf(`%d`, index),
				pp.Sprintf(`%s`, line),
			)
		}
	}
	if false {
		fmt.Println(
			pp.Sprintf("%s", c),
		)
	}
	if false {
		os.Exit(1)
	}
	return &c
}

func monitor_keystrokes(conn *itermctl.Connection) {
	keystrokes, err := itermctl.MonitorKeystrokes(context.Background(), conn, itermctl.AllSessions)
	if err != nil {
		panic(err)
	}

	for ks := range keystrokes {
		started := time.Now()
		chars := ks.GetCharacters()
		windowid, err := _app.ActiveTerminalWindowId()
		F(err)
		sess := _app.ActiveSession()
		sessid := sess.Id()
		qty_lines, err := sess.NumberOfLines()
		F(err)
		tabid, err := _app.SelectedTabId()
		F(err)
		contents, err := sess.ScreenContents(nil)
		F(err)
		lines := []string{}
		latest_lines := []string{}
		vim_match := false
		dev_match := false

		for index, line := range contents.GetContents() {
			if int(index) <= int(qty_lines.History) {
				line := fmt.Sprintf("%s", line.GetText())
				lines = append(lines, line)
				latest_lines = append(latest_lines, line)
			}
		}
		for _, line := range latest_lines {
			for _, m := range MATCHED_VIM_SWAP_STRINGS {
				if strings.Contains(line, `_app.ActiveTerminalWindowId()`) {
					dev_match = true
				}
				if strings.Contains(line, m) {
					vim_match = true
				}
			}
		}
		if dev_match {
			vim_match = false
		}
		vim_debug := ``
		if vim_match {
			RemoveEmpty(&latest_lines)
			vim_debug = fmt.Sprintf(`%s`, strings.Join(latest_lines, "\n"))

			vim_crash := parse_crash_msg(vim_debug)

			//			vims, err := get_remote_vms(VIM_LOCAL_PORT)
			//		F(err)
			iterm_profile, err := get_iterm_profile_from_vim_crash(vim_crash)
			if err != nil {
				continue

			}
			if false {
				fmt.Println(
					"vim crash:  ", pp.Sprintf(`%s`, vim_crash),
					//		"remote vims:  ", pp.Sprintf(`%s`, vims),
					"iterm profile:  ", pp.Sprintf(`%s`, iterm_profile),
				)
			}
			new_sess := _app.Session(iterm_profile.Session)

			if new_sess == nil {
				F(errors.New(fmt.Sprintf(`invalid session id %s`, iterm_profile.Session)))
			}

			new_contents, err := new_sess.ScreenContents(nil)
			F(err)
			new_lines := []string{}

			for _, _line := range new_contents.GetContents() {
				new_lines = append(new_lines, fmt.Sprintf("%s", _line.GetText()))
			}

			new_output := fmt.Sprintf("%s", strings.Join(new_lines, "\n"))

			msg := fmt.Sprintf(`


       Activating Existing Iterm2 Session

Vim Crash
	|   File:                     %s

Current Screen Contents:

%s

New Screen Contents:

%s

Current Session:
	|   Session:                  %s
	| # Window:                   %s
	| # Tab:                      %s

New Session:
	|   Session:                  %s
	| # Window:                   %d
	| # Tab:                      %d
  |


	`,
				pterm.FgLightMagenta.Sprintf(vim_crash.VimFilePath),

				pterm.FgLightRed.Sprintf(vim_debug),

				pterm.FgLightGreen.Sprintf(new_output),

				pterm.FgLightCyan.Sprintf(sessid),
				pterm.FgLightRed.Sprintf(windowid),
				pterm.FgLightGreen.Sprintf(tabid),

				pterm.FgLightCyan.Sprintf(iterm_profile.Session),
				iterm_profile.Window,
				iterm_profile.Tab,
			)
			if KEYSTROKE_ACTION_DEBUG_RESULT {
				pr(msg)
			}

			if KEYSTROKE_ACTION_SHOW_ALERT {
				button, err := _app.ShowAlert(itermctl.Alert{
					Title:    "VIM Crash",
					Subtitle: fmt.Sprintf("VIM Crash Detected: %s", `xxxxxxxxxx`),
				}, windowid)
				F(err)
				pp.Println(button)
			}
			if KEYSTROKE_ACTION_SELECT_NEW_SESSION {

				sel_new_session := _app.Session(iterm_profile.Session)
				err := sel_new_session.Activate()
				if err != nil {
					fmt.Println(fmt.Sprintf(`Error selecting new session %s: %s`, iterm_profile.Session, err))
					F(err)
				}

				sel_tab := fmt.Sprintf(`%d`, iterm_profile.Tab)
				if KEYSTROKE_ACTION_DEBUG_DETAILED {
					pp.Println(
						"Selecting tab", iterm_profile.Tab,
						"type: ", fmt.Sprintf(`%T`, iterm_profile.Tab),
						"sel_tab type: ", fmt.Sprintf(`%T`, sel_tab),
						"sel_tab val: ", fmt.Sprintf(`%v`, sel_tab),
					)
				}
				if KEYSTROKE_ACTION_SELECT_NEW_TAB {
					sel_err := _app.SelectTab(sel_tab)
					F(sel_err)
				}
				if KEYSTROKE_ACTION_QUIT_OLD_SESSION_VIM {
					sess.SendText(VIM_QUIT_TEXT, false)

				}

			}
		}
		if KEYSTROKE_ACTION_DEBUG_DETAILED {

			fmt.Printf(`

typed: %s
  | window id: %s
  | session id: %s
  | tab id: %s
  | keystroke parse dur: %s
  | Qty Lines:                     %d
  | Lines Qty:            %d
  | Latest Lines Qty:             %d
  | Latest Lines Bytes:           %d
  | Vim Case match?               %v

  %s

        `,
				chars,
				windowid,
				sessid,
				tabid,
				time.Since(started),
				qty_lines.History,
				len(lines),
				len(latest_lines),
				len(fmt.Sprintf(`%s`, latest_lines)),
				vim_match,
				err_msg(vim_debug),
			)

		}

	}
}
