package osxkeystrokes

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"dev.local/utils"
	"github.com/hpcloud/tail"
	"github.com/k0kubun/pp"
)

var KEYLOGGER_CMD = `./keylogger`

func dump_keystroke_log_to_stdout() {
	for {
		fmt.Println(pp.Sprintf(`%s`, dump_keystroke_log()))

		time.Sleep(5 * time.Second)
	}
}

func dump_keystroke_log() string {
	m.Lock()
	defer m.Unlock()
	r := ``
	pp.Println(keystrokes)
	return r
}

var m sync.Mutex

func add_ks(ks string) {
	m.Lock()
	defer m.Unlock()
	return
}

var keystrokes = []KeyStrokeLogEntry{}

func init() {
	for _, v := range active_keystroke_lengths {
		keystrokes = append(keystrokes, KeyStrokeLogEntry{Length: v, Keystrokes: []string{}})
	}
}

type KeyStrokeLogEntry struct {
	Length     int
	Keystrokes []string
}

var active_keystroke_lengths = []int{1, 2, 3, 4, 5}

func stderr_logger(msg string) {
	pp.Println("ERR>> ", msg)
}

func stdout_logger(msg string) {
	pp.Println("OUT>> ", msg)
}

func StderrLogger(stderr_logger func(string)) {
	go utils.ExecAsync(stdout_logger, stderr_logger, KEYLOGGER_CMD)
	for {
		time.Sleep(10 * time.Second)
	}
	//	dump_keystroke_log_to_stdout()

}

func Logger() {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "keylogger-")
	utils.F(err)
	tmp_name := tmpFile.Name()
	go utils.ExecAsync(stdout_logger, stderr_logger, `/usr/local/bin/keylogger`, tmp_name)
	go dump_keystroke_log_to_stdout()
	t, err := tail.TailFile(tmp_name, tail.Config{Follow: true})
	utils.F(err)
	for line := range t.Lines {
		fmt.Println(">> ", line.Text)
		if len(line.Text) > 0 {
			add_ks(line.Text)
		}
	}
}
