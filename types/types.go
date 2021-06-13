package types

import "time"

type ItermProfile struct {
	LocalPort        int    `json:"local_port";`
	RemotePort       int    `json:"remote_port";`
	User             string `json:"user";`
	RemoteListenHost string `json:"remote_listen_host";`
	Pid              uint   `json:"pid";`

	Cwd         string
	SwapFile    string
	VimFile     string
	VimFilePath string

	Window  int
	Session string
	Tab     int

	Hostname           string
	IP                 string
	VimRunningDuration time.Duration
}

type CurrentFocus struct {
	Session string
	Window  string
	Tab     string
}
