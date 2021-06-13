package main

import (
	"fmt"

	"dev.local/osxkeystrokes"
	"github.com/k0kubun/pp"
)

func key_logged(msg string) {
	pp.Println("Key Logged:-->     ", msg)
}

func main() {
	osxkeystrokes.StderrLogger(key_logged)

	fmt.Println(`OK`)
}
