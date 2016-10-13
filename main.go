package main

import "fmt"

var version = "v0.1.1"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("rip %s%s",
		version,
		dirty)
	Execute(displayVersion)
}
