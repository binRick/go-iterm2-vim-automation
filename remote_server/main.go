package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jackpal/gateway"
	"github.com/shirou/gopsutil/v3/process"
)

func main() {
	go webserver()
	for {
		find_vims()
		time.Sleep(10 * time.Second)
	}

}

func find_vims() (*[]ItermProfile, error) {
	procs, _ := process.Processes()
	iterm_profiles := []ItermProfile{}
	for _, proc := range procs {
		var ssh_client_port int64 = 0
		var raw_iterm_session_id string
		port_avail := false
		var iterm_session_id string
		var iterm_profile ItermProfile
		var iterm_profile_vars_encoded string
		var open_swap_file string

		n, _ := proc.Name()

		if n != `vim` {
			continue
		}

		pe, err := ReadProcessEnvironment(int64(proc.Pid))
		F(err)
		for _, e := range NullTermToStrings(pe) {
			if strings.HasPrefix(e, `ITERM_PROFILE_VARS_ENCODED=`) {
				enc_spl := strings.Split(e, `=`)
				enc_k := enc_spl[0] + `=`

				iterm_profile_vars_encoded = strings.Replace(e, enc_k, ``, 1)

				iterm_profile_vars_decoded, err := base64.StdEncoding.DecodeString(iterm_profile_vars_encoded)
				F(err)

				uerr := json.Unmarshal(iterm_profile_vars_decoded, &iterm_profile)
				F(uerr)
			}
			if strings.HasPrefix(e, `ITERM_SESSION_ID=`) {
				raw_iterm_session_id = strings.Split(e, `=`)[1]
			}
			if strings.HasPrefix(e, `SSH_CLIENT=`) {
				ssh_client_port, _ = strconv.ParseInt(strings.Split(strings.Split(e, `=`)[1], ` `)[1], 10, 0)
				F(err)
			}
		}
		if iterm_profile.Pid < 1 {

			continue
		}
		if ssh_client_port < 1 || raw_iterm_session_id == `` {
			continue
		}

		of, err := proc.OpenFiles()
		F(err)
		for _, f := range of {
			if strings.HasSuffix(f.Path, `.swp`) {
				open_swap_file = f.Path
			}
		}
		if open_swap_file == `` {
			continue
		}
		vim_file := fmt.Sprintf(`%s/%s`,
			path.Dir(open_swap_file),
			path.Base(open_swap_file)[1:len(path.Base(open_swap_file))-4],
		)

		info, err := os.Stat(vim_file)
		if os.IsNotExist(err) {
			continue
		}
		if info.IsDir() {
			continue
		}
		w_letters := ``
		t_letters := ``
		iterm_loc := strings.Split(raw_iterm_session_id, `:`)[0]
		iterm_session_id = strings.Split(raw_iterm_session_id, `:`)[1]
		var window_id int = -1
		var tab_id int = -1
		for _, c := range iterm_loc {
			_, err := strconv.ParseInt(string(c), 10, 0)
			is_number := false
			if err == nil {
				is_number = true
			} else {
				is_number = false
			}
			if is_number && len(w_letters) > 0 && tab_id == -1 {
				t_letters = fmt.Sprintf(`%s%s`, t_letters, string(c))
				continue
			}

			if is_number && window_id == -1 && len(t_letters) == 0 {
				w_letters = fmt.Sprintf(`%s%s`, w_letters, string(c))
				continue
			}
			if !is_number && len(w_letters) > 0 && len(t_letters) > 0 {
				t_int, t_err := strconv.ParseInt(string(t_letters), 10, 0)
				F(t_err)
				tab_id = int(t_int)
				w_int, t_err := strconv.ParseInt(string(w_letters), 10, 0)
				F(t_err)

				window_id = int(w_int)
				break
			}

		}
		proc_created, err := proc.CreateTime()
		F(err)
		proc_created_time := time.Unix(proc_created/1000, 0)
		proc_created_duration := time.Since(proc_created_time)

		con_host := net.JoinHostPort(iterm_profile.RemoteListenHost, fmt.Sprintf(`%d`, iterm_profile.RemotePort))
		gateway_interface, err := gateway.DiscoverInterface()
		F(err)
		hostname, err := os.Hostname()
		F(err)
		u, err := proc.Username()
		F(err)
		cwd, err := proc.Cwd()
		F(err)
		conn, err := net.DialTimeout("tcp", con_host, time.Millisecond*500)
		if err == nil {
			port_avail = true
			conn.Close()
		}
		iterm_profile.Window = window_id
		iterm_profile.Session = iterm_session_id
		iterm_profile.Tab = tab_id
		iterm_profile.Cwd = cwd
		iterm_profile.SwapFile = open_swap_file
		iterm_profile.VimFile = path.Base(vim_file)
		iterm_profile.VimFilePath = vim_file
		iterm_profile.Hostname = hostname
		iterm_profile.IP = fmt.Sprintf(`%s`, gateway_interface)
		iterm_profile.VimRunningDuration = proc_created_duration

		fmt.Println(fmt.Sprintf(`
 |> Process:
	| Program:                %s 
	| PID:                    %d 
 |> User:
	| User:                   %s
 |> Server:
	| Hostname:               %s
 |> Session:
	| Address:                %s
	| SSH Client Port:        %d 
	| Cwd:                    %s
 |> Vim:
	| Swap File Path:         %s
	| Vim File Path:          %s
	| Vim File:               %s
	| Age:										%s
 |> Iterm2:
	| Window ID:              %d
	| Tab ID:                 %d
	| Raw Iterm2 Session:     %s 
	| Iterm2 Session:         %s 
 |> Decoded Session:
	| Pid:                    %d
	| My Local Port:          %d
	| My Local Host:          %s
	| Remote  Port:           %d
	| Port Available?         %v


`,
			n,
			proc.Pid,
			u,
			hostname,
			gateway_interface,
			ssh_client_port,
			cwd,

			open_swap_file,
			vim_file,
			path.Base(vim_file),
			proc_created_duration,

			window_id,
			tab_id,
			raw_iterm_session_id,
			iterm_session_id,

			iterm_profile.Pid,
			iterm_profile.RemotePort,
			iterm_profile.RemoteListenHost,
			iterm_profile.LocalPort,
			port_avail,
		))
		iterm_profiles = append(iterm_profiles, iterm_profile)

	}

	return &iterm_profiles, nil
}
