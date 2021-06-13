package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func ExecAsync(_cmd string, stdout_callback func(string), stderr_callback func(string)) {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf(`%s`, _cmd))
	var wg sync.WaitGroup
	o := make(chan struct{})
	e := make(chan struct{})
	wg.Add(1)
	go func(cmd *exec.Cmd, c chan struct{}) {
		defer wg.Done()
		stderr, err := cmd.StderrPipe()
		F(err)
		<-c
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			m := scanner.Text()
			stderr_callback(m)
		}
	}(cmd, e)
	wg.Add(1)
	go func(cmd *exec.Cmd, c chan struct{}) {
		defer wg.Done()
		stdout, err := cmd.StdoutPipe()
		F(err)
		<-c
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			stdout_callback(m)
		}
	}(cmd, o)
	o <- struct{}{}
	e <- struct{}{}
	cmd.Start()
	wg.Wait()
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
