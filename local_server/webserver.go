package main

import (
	"encoding/json"
	"fmt"
	"net/http"

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
