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
	var add bool
	flag.BoolVar(&add, "add", false, "add new record")
	var pat string
	flag.StringVar(&pat, "find", "", "find record pattern")
	var cipher string
	flag.StringVar(&cipher, "cipher", "aes", "cipher algorithm AES, NaCI")
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbose level")
	flag.Parse()
	if version {
		fmt.Println(info())
		os.Exit(0)

	}

	// log time, filename, and line number
	if verbose > 0 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// determine vault location and if it is not provided or does not exists
	// creat $HOME/.pwm area and assign new vault file there
	_, err := os.Stat(vault)
	if vault == "" || os.IsNotExist(err) {
		udir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		vdir := fmt.Sprintf("%s/.pwm", udir)
		if verbose > 0 {
			log.Println("create vault dir", vdir)
		}
		err = os.MkdirAll(vdir, 0755)
		if err != nil {
			log.Fatal(err)
		}
		vault = fmt.Sprintf("%s/vault.aes", vdir)
	}

	// get vault secret
	salt, err := secret(verbose)
	if err != nil {
		log.Fatal(err)
	}

	// get vault records
	records, err := read(vault, salt, cipher, verbose)
	if err != nil {
		log.Fatal("unable to read vault, error ", err)
	}

	// perform vault operation
	if add {
		rec, err := input(verbose)
		if err != nil {
			log.Fatal(err)
		}
		newRecords := update(rec, records, verbose)
		write(vault, salt, cipher, newRecords, verbose)
		//         return
	}

	records, err = read(vault, salt, cipher, verbose)
	if err != nil {
		log.Fatal("unable to read vault, error ", err)
	}
	app := tview.NewApplication()
	//     listForm(app, records)
	gridView(app, records)

	//     find(vault, salt, cipher, pat, verbose)
}
