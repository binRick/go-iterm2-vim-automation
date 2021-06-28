package main

import (
	"sync"
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
type keySlice []string
