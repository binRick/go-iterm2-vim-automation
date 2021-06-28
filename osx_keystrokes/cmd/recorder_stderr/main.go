package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"time"

	"dev.local/osxkeystrokes"
	"dev.local/utils"
	"github.com/k0kubun/pp"
	"github.com/paulbellamy/ratecounter"
	"gopkg.in/yaml.v2"
)

var HELD_KEYS = []string{
	`left-option`,
	`right-option`,
	`left-shift`,
	`left-cmd`,
	`right-cmd`,
	`right-shift`,
	`left-ctrl`,
	`right-ctrl`,
	`caps`,
}

var key_codes KeyCodes

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

func (hks *HeldKeys) GetHeld() []string {
	held := []string{}
	for _, hk := range *hks {
		if hk.IsHeld() {
			held = append(held, hk.Key)
		}
	}
	return held
}

var held_keys = HeldKeys{}

func (hks *HeldKeys) KeyIsHeld(k string) bool {
	h := false
	for _, hk := range *hks {
		if hk.Key == fmt.Sprintf(`[%s]`, k) && hk.IsHeld() {
			h = true
		}
	}
	return h
}
func (hk *HeldKey) IsHeld() bool {
	found := false
	for _, HELD_KEY := range HELD_KEYS {
		if hk.Key == fmt.Sprintf(`[%s]`, HELD_KEY) {
			found = true
		}
	}
	if !found {
		return false
	}
	ih := (hk.Qty % 2) == 0
	if ih {
		hk.LastDepress = time.Now()
	}
	return ih
}

func parse_key_codes() {
	data, err := ioutil.ReadFile(`key_codes.yaml`)
	utils.F(err)
	uerr := yaml.Unmarshal([]byte(data), &key_codes)
	utils.F(uerr)
	//pp.Println(key_codes)

}

func key_logged(msg string) {
	rate_counter.Incr(1)
	//shift_is_held := false
	hash := fmt.Sprintf(`%x`, md5.Sum([]byte(fmt.Sprintf(`%s`, msg))))
	_, has := held_keys[hash]
	if !has {
		held_keys[hash] = &HeldKey{
			Key:  msg,
			Hash: hash,
			Qty:  1,
			Rate: ratecounter.NewRateCounter(1000 * time.Millisecond),
			//pp.Println(HV)
		}
	}
	hk, _ := held_keys[hash]
	hk.Qty = hk.Qty + 1
	hk.Rate.Incr(1)

	for k, _ := range rate_counters {
		if msg == fmt.Sprintf(`[%s]`, k) {
			rate_counters[k].Incr(1)
		}
	}
	fmt.Println(fmt.Sprintf(`
| =============================
| Key Logged:                 %s
| Key Rate/sec: v             %d
| == Modifiers
|  Left Shift?                %v
|  Left Control?              %v
|  Left Option?               %v
|  Left Command?              %v
|  Caps?                      %v
| ==
| Held Key:
|  Had?                       %v
|  Hash:                      %v
|  Key:                       %v
|  Qty:                       %v
|  Rate:                      %v
|  Time Since Last Depress:   %s
|  Is Held?                   %v
| ==
| Clear Key Rate/sec:        %d
| Equals Key Rate/sec:       %d
| Divide Key Rate/sec:       %d
| Asterisk Key Rate/sec:     %d
| =============================

	`,
		pp.Sprintf(`%s`, msg),
		rate_counter.Rate()*1,

		held_keys.KeyIsHeld(`left-shift`),
		held_keys.KeyIsHeld(`left-ctrl`),
		held_keys.KeyIsHeld(`left-option`),
		held_keys.KeyIsHeld(`left-cmd`),
		held_keys.KeyIsHeld(`caps`),

		has,
		hash,
		hk.Key,
		hk.Qty,
		hk.Rate.Rate(),
		time.Since(hk.LastDepress),
		hk.IsHeld(),

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
