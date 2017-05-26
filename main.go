package main

import "fmt"

var version = "v0.2.0"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("rip %s%s",
		version,
		dirty)
	Execute(displayVersion)
}
