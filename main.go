package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(1)
	}

	config := loadConfig(flag.Arg(0))
	site := newSite(config)
	site.build()
}

func usage() {
	fmt.Println("Usage: blgo <config_file>")
}
