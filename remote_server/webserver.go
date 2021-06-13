package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var WEBSERVER_PORT = 15223

func HandleListRequest(w http.ResponseWriter, r *http.Request) {
	vims, err := find_vims()
	F(err)

	if err := json.NewEncoder(w).Encode(vims); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	return
}
func HandleTestRequest(w http.ResponseWriter, r *http.Request) {
	response := `ok`
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, string(response))
}

func HandleVimFilePidRequest(w http.ResponseWriter, r *http.Request) {
	file_name := mux.Vars(r)["file_name"]
	response := `hello ` + string(file_name)
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, string(response))
}

func add_routes(router *mux.Router) {
	router.HandleFunc("/", HandleTestRequest).Methods(http.MethodGet)
	router.HandleFunc("/list", HandleListRequest).Methods(http.MethodGet)
	router.HandleFunc("/api/vim_file_pid/{file_name:.*}", HandleVimFilePidRequest).Methods(http.MethodGet)
}

func webserver() {
	router := mux.NewRouter()
	router.SkipClean(true)
	add_routes(router)
	addr := fmt.Sprintf(`0.0.0.0:%d`, WEBSERVER_PORT)
	fmt.Println(fmt.Sprintf("Listening on %s", addr))
	http.ListenAndServe(addr, router)
}
