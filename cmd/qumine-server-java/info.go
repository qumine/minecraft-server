package main

import (
	"flag"
	"fmt"
	"os"
)

func showUsage() {
	flag.Usage()
	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%v, commit %v, built at %v", version, commit, date)
	os.Exit(0)
}
