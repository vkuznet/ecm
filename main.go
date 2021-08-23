package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// version of the code
var gitVersion string

// Info function returns version string of the server
func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("pwm git=%s go=%s date=%s", gitVersion, goVersion, tstamp)
}

func main() {
	var vault string
	flag.StringVar(&vault, "vault", "", "vault name")
	var secret string
	flag.StringVar(&secret, "secret", "", "vault secret")
	var add bool
	flag.BoolVar(&add, "add", false, "add new record")
	var pat string
	flag.StringVar(&pat, "find", "", "find record pattern")
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	flag.Parse()
	if version {
		fmt.Println(info())
		os.Exit(0)

	}

	if vault == "" {
		log.Fatal("empty vault")
	}
	if add {
		write(vault, secret)
		return
	}
	find(vault, secret, pat)
}
