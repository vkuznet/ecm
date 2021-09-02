package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/rivo/tview"
)

// version of the code
var gitVersion, gitTag string

// Info function returns version string of the server
func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("pwm git=%s tag=%s go=%s date=%s", gitVersion, gitTag, goVersion, tstamp)
}

func main() {
	var vname string
	flag.StringVar(&vname, "vault", "", "vault name")
	var add bool
	flag.BoolVar(&add, "add", false, "add new record")
	var pat string
	flag.StringVar(&pat, "find", "", "find record pattern")
	var cipher string
	flag.StringVar(&cipher, "cipher", "", "cipher to use (aes, nacl)")
	var dfile string
	flag.StringVar(&dfile, "decrypt", "", "decrypt given file name")
	var efile string
	flag.StringVar(&efile, "encrypt", "", "encrypt given file and place it into vault")
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbose level")
	flag.Parse()
	if version {
		fmt.Println(info())
		os.Exit(0)

	}

	// decrypt record
	if dfile != "" {
		decryptFile(dfile, cipher)
		os.Exit(0)
	}

	// parse input config
	configFile := fmt.Sprintf("%s/config.json", pwmHome())
	err := ParseConfig(configFile, verbose)
	if err != nil {
		log.Fatal(err)
	}

	// set Theme for our app
	setTheme("grey")

	// get vault secret
	salt, err := secret(verbose)
	if err != nil {
		log.Fatal(err)
	}

	// initialize our vault
	vault := Vault{Cipher: cipher, Secret: salt, Verbose: verbose}

	// create our vault
	err = vault.Create(vname)
	if err != nil {
		log.Fatalf("unable to create vault, error %v", err)
	}

	// encrypt given record
	if efile != "" {
		vault.EncryptFile(efile)
	}

	// read from our vault
	err = vault.Read()
	if err != nil {
		log.Fatal("unable to read vault, error ", err)
	}

	// create vault app and run it
	app := tview.NewApplication()
	gridView(app, &vault)
}
