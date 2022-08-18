package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

// version of the code
var gitVersion, gitTag string

// ecmInfo function returns version string of the server
func ecmInfo() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("ecm git=%s tag=%s go=%s date=%s", gitVersion, gitTag, goVersion, tstamp)
}

func main() {
	var version bool
	flag.BoolVar(&version, "version", false, "show version")
	var config string
	flag.StringVar(&config, "config", "", "start HTTP server with provided configuration")
	flag.Parse()
	if version {
		fmt.Println(ecmInfo())
		os.Exit(0)

	}

	// start server
	startServer(config)
}
