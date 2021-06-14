package main

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"dev.local/osxkeystrokes"
	"dev.local/utils"
	"github.com/k0kubun/pp"
	"gopkg.in/yaml.v2"
)

type KeyCode struct {
	Name   string   `yaml:"Name";`
	Action string   `yaml:"Action";`
	Keys   keySlice `yaml:"Keys";`
}

type LastKeyTracker struct {
	LastKeys  keySlice
	MaxLength int
	m         sync.Mutex
}
type KeyTracker LastKeyTracker
type KeyCodes *[]KeyCode

var key_codes KeyCodes

type keySlice []string

var key_tracker = KeyTracker{MaxLength: 10}

func monitor_tracked_keys() {
	for {
		key_tracker.m.Lock()
		pp.Println(key_tracker)
		key_tracker.m.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func init() {
	//go monitor_tracked_keys()
}

func parse_key_codes() {
	data, err := ioutil.ReadFile(`key_codes.yaml`)
	utils.F(err)
	uerr := yaml.Unmarshal([]byte(data), &key_codes)
	utils.F(uerr)
	pp.Println(key_codes)

}

func key_logged(msg string) {
	pp.Println("Key Logged:-->     ", msg)
}

func main() {
	osxkeystrokes.StderrLogger(key_logged)

	fmt.Println(`OK`)
}
