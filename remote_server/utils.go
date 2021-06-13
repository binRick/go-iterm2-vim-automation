package main

import (
	"bytes"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

func ReadProcessEnvironment(pid int64) ([]byte, error) {
	proc_path := fmt.Sprintf(`/proc/%d/environ`, pid)
	b, err := ioutil.ReadFile(proc_path)
	if err != nil {
		return b, err
	}
	return b, nil
}

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

func F(err error) {
	if err != nil {
		log.Error(err)
		panic(err)
	}
}
