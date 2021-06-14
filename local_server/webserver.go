package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/k0kubun/pp"
	log "github.com/sirupsen/logrus"
)

var WEBSERVER_PORT = 48923

type Iterm2 struct {
	Window       *string
	WindowNumber *int32
	Session      string
	Tab          *string
}

func ListIterm2() *[]Iterm2 {
	r := []Iterm2{}

	cur_sessions, err := _app.ListSessions()
	F(err)
	for _, w := range cur_sessions.Windows {
		if false {
			fmt.Println(
				fmt.Sprintf(`%T`, w),
				pp.Sprintf(`%s`, w),
			)
		}

		for _, t := range w.Tabs {

			sess := t.Root.Links[0].Child
			if false {
				fmt.Println(
					fmt.Sprintf(`%T`, t.Root.Links[0].Child),
					//pp.Sprintf(`%s`, t.Root.Links[0].Child),
					pp.Sprintf(`%s`, sess),
				)
			}
			r = append(r, Iterm2{
				Window:       w.WindowId,
				Session:      `yyyyyyyyy`,
				Tab:          t.TabId,
				WindowNumber: w.Number,
			})
		}
	}

	return &r
}

type Iterm2NewTabResponse struct {
	Hostname      string
	Directory     string
	Cmd           string
	NewTabIndex   uint32
	WindowID      string
	NewTabProfile string
	SessionID     string

	Result   string
	SendText string
}

func HandleNewIterm2TabRequest(w http.ResponseWriter, r *http.Request) {

	_q_hostname, _q_hostname_ok := r.URL.Query()["hostname"]
	_q_directory, _q_directory_ok := r.URL.Query()["directory"]
	_q_cmd, _q_cmd_ok := r.URL.Query()["cmd"]

	pp.Println(
		_q_hostname, _q_hostname_ok,
		_q_directory, _q_directory_ok,
		_q_cmd, _q_cmd_ok,
	)

	_cmd, err := url.QueryUnescape(_q_cmd[0])
	F(err)
	cmd := _cmd

	_hostname, err := url.QueryUnescape(_q_hostname[0])
	F(err)
	hostname := _hostname

	_directory, err := url.QueryUnescape(_q_directory[0])
	F(err)
	directory := _directory
	pp.Println(
		_directory,
		directory,
	)
	windowid, err := _app.ActiveTerminalWindowId()
	F(err)
	sess := _app.ActiveSession()
	sessid := sess.Id()
	cur_tab, err := _app.SelectedTabId()
	F(err)
	pp.Println(cur_tab)
	var new_tab_index uint32 = 1
	///	pp.Println(sess)
	res := Iterm2NewTabResponse{
		Hostname:      hostname,
		Directory:     directory,
		Cmd:           cmd,
		NewTabProfile: `Goonies`,
		NewTabIndex:   new_tab_index,
		WindowID:      windowid,
		SessionID:     sessid,
	}
	pp.Println(res)
	new_tab_response, err := _app.CreateTab(res.WindowID, res.NewTabIndex, res.NewTabProfile)
	F(err)
	//	pp.Println(new_tab_response)

	res.Result = fmt.Sprintf("Created Tab #%d", new_tab_response.TabId)

	if res.Cmd == `` {
		res.Cmd = `echo OK`
	}
	if res.Hostname == `localhost` {

	} else {
	}
	if len(res.Directory) > 0 {
		res.SendText = strings.Replace(fmt.Sprintf(`cd %s && %s`, strings.Replace(strings.Replace(res.Directory, "\n", "", -1), "\r", "", -1), res.Cmd), "\n", "", -1)
	}
	pp.Println(res.SendText)

	if len(res.SendText) > 0 {
		_app.ActiveSession().SendText(res.SendText+"\n", false)
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}
func HandleListVims(w http.ResponseWriter, r *http.Request) {
	vims, err := get_remote_vms(VIM_LOCAL_PORT)
	F(err)

	if err := json.NewEncoder(w).Encode(vims); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}

func HandleListIterm2(w http.ResponseWriter, r *http.Request) {

	if err := json.NewEncoder(w).Encode(ListIterm2()); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return

}
func add_routes(router *mux.Router) {
	router.HandleFunc("/api/chrome/new_tab/{url:.*}", HandleWebserverRequest).Methods(http.MethodGet)
	router.HandleFunc("/api/iterm2/new_tab", HandleNewIterm2TabRequest).Methods(http.MethodGet)
	router.HandleFunc("/api/vim/test", HandleWebserverVimTest).Methods(http.MethodGet)
	router.HandleFunc("/api/vims/list", HandleListVims).Methods(http.MethodGet)
	router.HandleFunc("/api/iterm2/list", HandleListIterm2).Methods(http.MethodGet)
	router.HandleFunc("/api/iterm2/activate/window/{window_id:.*}/session/{session_id:.*}/tab/{tab_id:.*}", HandleActivateWindowID).Methods(http.MethodGet)
	router.HandleFunc("/api/iterm2/test/window/{window_id:.*}/session/{session_id:.*}/tab/{tab_id:.*}", HandleTestWindowID).Methods(http.MethodGet)
	router.HandleFunc("/api/iterm2/close/window/{window_id:.*}/session/{session_id:.*}", HandleCloseWindowID).Methods(http.MethodGet)

}
func webserver() {
	router := mux.NewRouter()
	router.SkipClean(true)
	add_routes(router)

	addr := fmt.Sprintf(`0.0.0.0:%d`, WEBSERVER_PORT)
	fmt.Println(fmt.Sprintf("Listening on %s", addr))
	http.ListenAndServe(addr, router)
}
