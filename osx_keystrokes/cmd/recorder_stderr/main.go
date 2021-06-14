package main

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"dev.local/osxkeystrokes"
	"dev.local/utils"
	"github.com/k0kubun/pp"
	"github.com/paulbellamy/ratecounter"
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
		//pp.Println(key_tracker)
		key_tracker.m.Unlock()
		time.Sleep(5 * time.Second)
	}
}

func init() {
	//go monitor_tracked_keys()
}

var monitored_key_seq_names = []string{`clear`, `equals`, `divide`, `asterisk`}

var rate_counter = ratecounter.NewRateCounter(1000 * time.Millisecond)
var rate_counters = map[string]*ratecounter.RateCounter{
	`clear`:    ratecounter.NewRateCounter(1000 * time.Millisecond),
	`equals`:   ratecounter.NewRateCounter(1000 * time.Millisecond),
	`divide`:   ratecounter.NewRateCounter(1000 * time.Millisecond),
	`asterisk`: ratecounter.NewRateCounter(1000 * time.Millisecond),
}

//counter.Incr(1)
//counter.Rate()

func parse_key_codes() {
	data, err := ioutil.ReadFile(`key_codes.yaml`)
	utils.F(err)
	uerr := yaml.Unmarshal([]byte(data), &key_codes)
	utils.F(uerr)
	//pp.Println(key_codes)

}

func key_logged(msg string) {
	rate_counter.Incr(1)
	for k, _ := range rate_counters {
		if msg == fmt.Sprintf(`[%s]`, k) {
			rate_counters[k].Incr(1)
		}
	}
	fmt.Println(fmt.Sprintf(`
	| =============================
	| Key Logged:                %s
	| Key Rate/sec:              %d
	| ==
	| Clear Key Rate/sec:        %d
	| Equals Key Rate/sec:       %d
	| Divide Key Rate/sec:       %d
	| Asterisk Key Rate/sec:     %d
	| =============================

	`,
		pp.Sprintf(`%s`, msg),
		rate_counter.Rate()*1,

		rate_counters[`clear`].Rate(),
		rate_counters[`equals`].Rate(),
		rate_counters[`divide`].Rate(),
		rate_counters[`asterisk`].Rate(),
	))
}

func main() {
	osxkeystrokes.StderrLogger(key_logged)

	fmt.Println(`OK`)
}
