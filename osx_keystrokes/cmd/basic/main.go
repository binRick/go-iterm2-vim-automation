package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"dev.local/utils"
	"github.com/hpcloud/tail"
	"github.com/k0kubun/pp"
)

func stderr_logger(msg string) {
	pp.Println("ERR>> ", msg)
}

func stdout_logger(msg string) {
	pp.Println("OUT>> ", msg)
}

func main() {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "keylogger-")
	utils.F(err)
	tmp_name := tmpFile.Name()

	//utils.ExecAsync(osxkeystrokes.KEYLOGGER_CMD, stdout_logger, stderr_logger)
	utils.ExecAsync(stdout_logger, stderr_logger, `ls`)
	utils.ExecAsync(stdout_logger, stderr_logger, `ls`, `/`)
	utils.ExecAsync(stdout_logger, stderr_logger, `ls`, `/etc`)
	utils.ExecAsync(stdout_logger, stderr_logger, `ls`, `/etc1`)
	go utils.ExecAsync(stdout_logger, stderr_logger, `/usr/local/bin/keylogger`, tmp_name)
	t, err := tail.TailFile(tmp_name, tail.Config{Follow: true})
	utils.F(err)
	for line := range t.Lines {
		fmt.Println(">> ", line.Text)
	}
	//	utils.ExecAsync(`ls /1`, stdout_logger, stderr_logger)
}
