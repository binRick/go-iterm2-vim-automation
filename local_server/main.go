package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/k0kubun/pp"
	"github.com/mgutz/ansi"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"

	"github.com/andybrewer/mack"
	"github.com/shirou/gopsutil/process"
	"mrz.io/itermctl"
	"mrz.io/itermctl/rpc"
)

func init_mack() {
  	if false {
	}
	if true {
		mack.Notify("Complete")
		response, err := mack.Dialog("Enter a ToDo", "ToDo Wizard", "My new ToDo")
		F(err)
		pp.Println(response)
	}
	if true {
		list := mack.ListOptions{
			Items:         []string{"item one", "item two"},
			Title:         "My List Title",
			Message:       "Pick one or more items from this list",
			DefaultItems:  []string{"item one"},
			AllowMultiple: true,
		}
		selected, didCancel, err := mack.ListWithOpts(list)
		F(err)
		pp.Println(selected, didCancel)
	}
}
func init() {
if false {
	init_mack()
}
	go get_procs()
	tempFile, err := ioutil.TempFile("", "*-custom_control_test")
	F(err)

	_, err = tempFile.Write([]byte(ctrl_seq1.Escape("test-seq")))

	F(err)
	pp.Println(tempFile.Name())
}

func monitor_control_seq() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	re := regexp.MustCompile("^test-seq:.*")
	fmt.Println(
		"CONTROL_SEQUENCE_NAME:", CONTROL_SEQUENCE_NAME,
	)
	notifications, err := itermctl.MonitorCustomControlSequences(ctx, _conn, CONTROL_SEQUENCE_NAME, re, itermctl.AllSessions)
	F(err)

	msg := fmt.Sprintf(`
 ** Waiting for control sequeunces **

CONTROL_SEQUENCE_NAME: %v

`,
		CONTROL_SEQUENCE_NAME,
	)
	pp.Println(msg)
	dm := func() {
		select {
		case notification := <-notifications:
			pp.Println("New Sequence!")
			pp.Println(fmt.Sprintf(`Session: %s`, notification.Notification.GetSession()))
			pp.Println(fmt.Sprintf(`Matches qty: %d`, len(notification.Matches)))
			for _, m := range notification.Matches {
				pp.Println(fmt.Sprintf(`   Match:     %d bytes`, len(m)))
				if len(strings.Split(m, `:`)) == 2 {
					seq_json := map[string]interface{}{}
					seq_enc := strings.Split(m, `:`)[1]

					seq_dec, err := base64.StdEncoding.DecodeString(seq_enc)
					F(err)

					seq_dec_trimmed := strings.Trim(fmt.Sprintf(`%s`, seq_dec), "")

					dec_err := json.Unmarshal([]byte(seq_dec_trimmed), &seq_json)
					F(dec_err)

					pp.Println(
						seq_json,
					)
				}

			}
		}
	}

	for {
		dm()
	}
}

func (V *ActiveVim) IsValid() bool {
	v := false
	V.expires_ts = int64(V.ts) + int64(V.interval)
	if time.Now().Unix() < V.expires_ts {
		v = true
	}

	return v
}

type ActiveVim struct {
	pid        int64
	file       string
	ts         int64
	interval   int64
	expires_ts int64
}
type ActiveVims struct {
	Vims *ActiveVim
}

const (
	DEFAULT_FORCE_CLOSE = false
)

var (
	_conn         *itermctl.Connection
	_app          *itermctl.App
	ctrl_seq1     = itermctl.NewCustomControlSequenceEscaper(CONTROL_SEQUENCE_NAME)
	current_focus = &CurrentFocus{}
	err_msg       = ansi.ColorFunc("red")
)

type CurrentFocus struct {
	Session string
	Window  string
	Tab     string
}

var Clock = rpc.StatusBarComponent{
	ShortDescription: "VIM Crash Detector",
	Description:      "VIM Crash Detector",
	Exemplar:         "exemplar",
	UpdateCadence:    1,
	Identifier:       "vim.manager",
	Knobs: ClockKnobs{
		Location: "UTC",
		Option1:  `xxxxxxx`,
	},
	OnClick: OnClick,
	RPC: rpc.RPC{
		Name:     "itermctl_example_clock",
		Function: UpdateClock,
	},
}

type ClockKnobs struct {
	Location string
	Option1  string
}

var header = pterm.HeaderPrinter{
	TextStyle: pterm.NewStyle(pterm.FgYellow, pterm.Bold),
	Margin:    5,
}
var pr = header.Println

func NullTermToStrings(b []byte) (s []string) {
	nt := 0
	ntb := byte(nt)
	for {
		i := bytes.IndexByte(b, ntb)
		if i == -1 {
			break
		}
		s = append(s, string(b[0:i]))
		b = b[i+1:]
	}
	return
}

type SSHConnection struct {
	Pid        int64
	RemotePort uint
	RemoteHost string
	LocalPort  uint
}

var ssh_connections = []SSHConnection{}

func get_procs() {
	procs, _ := process.Processes()
	for _, proc := range procs {
		n, err := proc.Name()
		if err != nil {
			continue
		}

		if n != `ssh` {
			continue
		}

		conns, err := proc.Connections()
		F(err)

		for _, conn := range conns {
			if conn.Status == `ESTABLISHED` {
				ssh_connections = append(ssh_connections, SSHConnection{
					Pid:        int64(proc.Pid),
					LocalPort:  uint(conn.Laddr.Port),
					RemotePort: uint(conn.Raddr.Port),
					RemoteHost: fmt.Sprintf(`%s`, conn.Raddr.IP),
				})
			}
		}
	}
	for _, _ = range ssh_connections {
		//pp.Println(fmt.Sprintf(`[PID %d] :%d => %s:%d`, c.Pid, c.LocalPort, c.RemoteHost, c.RemotePort))
	}
}

func UpdateClock(invocation *rpc.Invocation) (interface{}, error) {
	knobs := &ClockKnobs{}
	err := invocation.Knobs(knobs)
	if err != nil {
		return nil, err
	}

	//location, err := time.LoadLocation(knobs.Location)

	//now := time.Now().In(location)
	return fmt.Sprintf("%s", `wow 123`), nil
	//return fmt.Sprintf("%s", now.Round(1*time.Second)), nil
}

func OnClick(invocation *rpc.Invocation) (interface{}, error) {
	args := rpc.ClickArgs{}
	if err := invocation.Args(&args); err != nil {
		return nil, fmt.Errorf("click handler: %w", err)
	}

	html := fmt.Sprintf("<p>WOW: %s</p>", args.SessionId)

	if err := invocation.OpenPopover(html, 320, 240); err != nil {
		return nil, fmt.Errorf("click handler: %w", err)
	}

	return nil, nil
}
func F(err error) {
	if err != nil {
		log.Error(err)
		panic(err)
	}
}

func exec_cmd(cmd string) (string, string, syscall.WaitStatus, error) {
	Cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf(`%s`, cmd))
	var stdout, stderr bytes.Buffer
	var waitStatus syscall.WaitStatus
	Cmd.Stdout = &stdout
	Cmd.Stderr = &stderr
	defer Cmd.Wait()
	if err := Cmd.Run(); err != nil {
		if err != nil {
			return string(stdout.Bytes()), string(stderr.Bytes()), waitStatus, err
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			return string(stdout.Bytes()), string(stderr.Bytes()), waitStatus, err
		}
	} else {
		waitStatus = Cmd.ProcessState.Sys().(syscall.WaitStatus)
		return string(stdout.Bytes()), string(stderr.Bytes()), waitStatus, err
	}

	return string(stdout.Bytes()), string(stderr.Bytes()), waitStatus, nil
}
func (cf *CurrentFocus) IsReady() bool {
	return len(cf.Session) > 16 && len(cf.Window) > 16 && len(cf.Tab) > 0
}

func monitor_focus() {
	notifications, err := itermctl.MonitorFocus(context.Background(), _conn)
	F(err)
	go func() {
		for {
			if current_focus.IsReady() {
				//pp.Println("current focus:   ", current_focus)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		if current_focus.Window == `` {
			windowId, err := _app.ActiveTerminalWindowId()
			F(err)
			current_focus.Window = string(windowId)
		}
		for notification := range notifications {
			if fmt.Sprintf(`%s`, notification.Which) == `WindowBecameKey` {
				current_focus.Window = string(notification.Id)
			}
			if fmt.Sprintf(`%s`, notification.Which) == `TabSelected` {
				current_focus.Tab = string(notification.Id)
			}
			if fmt.Sprintf(`%s`, notification.Which) == `SessionSelected` {
				current_focus.Session = string(notification.Id)
			}
		}
	}()

}

func HandleNumberOfLines(w http.ResponseWriter, r *http.Request) {
	session_id := mux.Vars(r)["session_id"]
	session := _app.Session(session_id)
	if session == nil {
		fmt.Println(fmt.Sprintf(`invalid session id %s`, session_id))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	nol, err := session.NumberOfLines()
	F(err)
	qty := fmt.Sprintf("%d\n", nol.History)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, qty)
	return

}

func HandleDumpSessionContents(w http.ResponseWriter, r *http.Request) {
	session_id := mux.Vars(r)["session_id"]
	session := _app.Session(session_id)
	if session == nil {
		fmt.Println(fmt.Sprintf(`invalid session id %s`, session_id))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sess_id := session.Id()

	contents, err := session.ScreenContents(nil)
	F(err)
	lines := []string{}

	for _, line := range contents.GetContents() {
		lines = append(lines, line.GetText())
	}
	pp.Println(session_id, sess_id, "contents lines qty:", len(lines))
	output := strings.Join(lines, "\n")
	output_enc := base64.StdEncoding.EncodeToString([]byte(output))

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, output_enc)
	return

}

func HandleTestWindowID(w http.ResponseWriter, r *http.Request) {
	window_id := mux.Vars(r)["window_id"]
	session_id := mux.Vars(r)["session_id"]
	tab_id := mux.Vars(r)["tab_id"]
	session := _app.Session(session_id)
	if session == nil {
		fmt.Println(fmt.Sprintf(`invalid session id %s`, session_id))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sess_id := session.Id()

	if false {
		sel_err := _app.SelectTab(tab_id)
		if sel_err != nil {
			fmt.Println(fmt.Sprintf(`invalid tab id %s: %s`, tab_id, sel_err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	contents, err := session.ScreenContents(nil)
	F(err)
	lines := []string{}

	for _, line := range contents.GetContents() {
		//lines = append(lines, fmt.Sprintf("%s", line.GetText()))
		lines = append(lines, line.GetText())
	}

	//output := fmt.Sprintf("%s", strings.Join(lines, "\n"))
	output := strings.Join(lines, "\n")
	output_enc := base64.StdEncoding.EncodeToString([]byte(output))

	if false {
		pp.Println(session, tab_id, session_id, window_id, sess_id, "lines qty:", len(lines))
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, output_enc)
	return

}

func RemoveEmpty(slice *[]string) {
	i := 0
	p := *slice
	for _, entry := range p {
		if strings.Trim(entry, " ") != "" {
			p[i] = entry
			i++
		}
	}
	*slice = p[0:i]
}

func HandleCloseWindowID(w http.ResponseWriter, r *http.Request) {
	window_id := mux.Vars(r)["window_id"]
	session_id := mux.Vars(r)["session_id"]
	if false {
		pp.Println(window_id, session_id)
	}
	force := DEFAULT_FORCE_CLOSE
	err := _app.CloseTerminalWindow(force, window_id)
	if err != nil {
		fmt.Println(fmt.Sprintf(`CloseTerminalWindow error: %s`, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := `OK`

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(response))
	return

}

/*
var (
	OPEN_SESSION_ACTIVATE   = true
	OPEN_SESSION_SELECT_TAB = true
)
*/
func HandleOpenSessionID(w http.ResponseWriter, r *http.Request) {
	session_id := mux.Vars(r)["session_id"]
	tab_id := mux.Vars(r)["tab_id"]

	msg := fmt.Sprintf(`Selecting Session %s | Tab %s | `, session_id, tab_id)

	//pp.Println(msg)
	session := _app.Session(session_id)
	pre_sess_active := session.Active()
	if !pre_sess_active { //&& OPEN_SESSION_ACTIVATE {
		F(session.Activate())
	}
	post_sess_active := session.Active()
	selected_tab, err := _app.SelectedTabId()
	F(err)

	if tab_id != selected_tab { //&& OPEN_SESSION_SELECT_TAB {
		F(_app.SelectTab(tab_id))
	}

	post_selected_tab, err := _app.SelectedTabId()
	F(err)

	fmt.Println(
		pp.Sprintf(`Pre session active? %v | `, pre_sess_active),
		pp.Sprintf(`Post session active? %v | `, post_sess_active),
		pp.Sprintf(`App Selected Tab: %s | `, selected_tab),
		pp.Sprintf(`Post App Selected Tab: %s | `, post_selected_tab),
		pp.Sprintf(`Requested Tab: %s | `, tab_id),
		//pp.Sprintf(`session: %s`, session),
	)
	response := fmt.Sprintf(`OK- %s`, msg)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(response))
	return
}
func HandleActivateWindowID(w http.ResponseWriter, r *http.Request) {
	window_id := mux.Vars(r)["window_id"]
	session_id := mux.Vars(r)["session_id"]
	tab_id := mux.Vars(r)["tab_id"]
	session := _app.Session(session_id)
	if false {
		pp.Println(window_id, session_id, tab_id)
	}
	err := session.Activate()
	if err != nil {
		fmt.Println(fmt.Sprintf(`%s`, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	F(_app.SelectTab(tab_id))
	F(session.SendText("vim xxxxxx", false))
	contents, err := session.ScreenContents(nil)

	F(err)
	for i, line := range contents.GetContents() {
		fmt.Printf("#%d >>> %s\n", i, line.GetText())
	}

	qty_lines, err := session.NumberOfLines()
	if err == nil {
		fmt.Println(fmt.Sprintf(`
Number Of Lines: %d

`, qty_lines.History))
	}

	sel_text, err := session.SelectedText()
	if err == nil {
		fmt.Println(fmt.Sprintf(`
Selected Text: %s

`, sel_text))
	}

	///escaped_str := ctrl_seq1.Escape(fmt.Sprintf("hello %s", "world"))
	escaped_str := ctrl_seq1.Escape("hello %s", "world")
	if false {
		pp.Println(`escaped str:    `, escaped_str)
	}

	response := `ok`
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(response))
	return
}

func HandleWebserverVimTest(w http.ResponseWriter, r *http.Request) {
	cur_sessions, err := _app.ListSessions()
	if false {
		pp.Println(fmt.Sprintf("%d Current Windows", len(cur_sessions.Windows)))
	}
	msg := fmt.Sprintf(`
    # Windows:                      %d
`,
		len(cur_sessions.Windows),
	)
	fmt.Println(msg)
	for _, w := range cur_sessions.Windows {
		msg := fmt.Sprintf(`       Window ID: %s
                               # Tabs:                                   %d
`,
			*w.WindowId,
			len(w.Tabs),
		)
		fmt.Println(msg)

	}

	windowId, err := _app.ActiveTerminalWindowId()
	F(err)
	if false {
		t, err := _app.GetText(itermctl.TextInputAlert{
			Title:        "Type something",
			Subtitle:     "Type something in the field below",
			Placeholder:  "Placeholder for your text",
			DefaultValue: "",
		}, windowId)
		F(err)
		if false {
			pp.Println(t)
		}
	}
	if false {
		button, err := _app.ShowAlert(itermctl.Alert{
			Title:    "Test",
			Subtitle: fmt.Sprintf("You typed: %s", `xxxxxxxxxx`),
		}, windowId)
		F(err)

		if false {
			pp.Println(button)
		}
	}
	if false {
		pp.Println(windowId)
		pp.Println(current_focus)
	}

	if err := json.NewEncoder(w).Encode(current_focus); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return

}

func HandleWebserverRequest(w http.ResponseWriter, r *http.Request) {
	//pp.Println("Request!")
	u := mux.Vars(r)["url"]
	q, err := url.QueryUnescape(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		F(err)
		return
	}
	//pp.Println(q)

	p, err := url.Parse(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		F(err)
		return
	}
	cmd := fmt.Sprintf(`chrome-cli open %s -n`, q)
	if p.Host != `` && p.Scheme != `` {
		_, _, _, err = exec_cmd(cmd)
		F(err)
	}
	//pp.Println(p)
	response := `hello`

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, string(response))

	fmt.Println(fmt.Sprintf("Responded to request for url %s with %d byte response", u, len(response)))

}

func main() {
	go webserver()

	conn, err := itermctl.GetCredentialsAndConnect("itermctl_statusbar_example", true)
	F(err)
	_conn = conn
	app, err := itermctl.NewApp(_conn)
	F(err)
	_app = app
	monitor_focus()
	if MONITOR_CONTROL_SEQUENCE {
		go monitor_control_seq()
	}
	monitor_keystrokes(conn)
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals
		conn.Close()
	}()

	err = rpc.RegisterStatusBarComponent(context.Background(), conn, Clock)
	if err != nil {
		panic(err)
	}

	conn.Wait()
}
