package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	//     "github.com/rivo/tview"

	vt "github.com/vkuznet/ecm/vault"
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
	var vname string
	flag.StringVar(&vname, "vault", "", "vault name")
	var cipher string
	flag.StringVar(&cipher, "cipher", "", "cipher to use (aes, nacl)")
	var version bool
	flag.BoolVar(&version, "version", false, "show version")
	var lockInterval int
	flag.IntVar(&lockInterval, "lock", 60, "lock interval in seconds")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbose level")
	flag.Parse()
	if version {
		fmt.Println(ecmInfo())
		os.Exit(0)

	}

	// use file name in a log
	if verbose > 0 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	// initialize our vault
	vault := vt.Vault{Cipher: getCipher(cipher), Verbose: verbose, Start: time.Now()}

	// create vault if necessary
	err := vault.Create(vname)
	if err != nil {
		log.Fatalf("unable to create vault, error %v", err)
	}

	// start term UI mode
	setTheme("grey")
	gpgApp(&vault, lockInterval)
}
