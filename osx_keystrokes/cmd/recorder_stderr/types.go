package main

import (
	"sync"
"time"
 "github.com/paulbellamy/ratecounter"

)

type HeldKey struct {
  Key         string
  Hash        string
  Qty         int64
  Rate        *ratecounter.RateCounter
  LastDepress time.Time
}

type HeldKeys map[string]*HeldKey


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
type keySlice []string
