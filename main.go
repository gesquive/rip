package main

import "os"
import "fmt"
import "path/filepath"

var version = "0.1.0"
var dirty = ""

func main() {
	displayVersion := fmt.Sprintf("%s v%s%s",
		filepath.Base(os.Args[0]),
		version,
		dirty)
	Execute(displayVersion)
}
