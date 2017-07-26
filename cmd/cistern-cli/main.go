package main

import (
	"flag"
	"fmt"
	"os"
)

var version = "0.1.0"

func main() {
	address := flag.String("address", "localhost:2020", "Cistern node address")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	_ = address
}
