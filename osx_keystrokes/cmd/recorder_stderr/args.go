package main

import (
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	name    = kingpin.Arg("name", "Name of user.").Required().String()
)

func init() {
	kingpin.Parse()
	fmt.Printf("%v, %s\n", *verbose, *name)
}
