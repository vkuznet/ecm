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
	//     log.Println("### existing records")
	//     for _, rec := range records {
	//         log.Println("rec", rec)
	//     }

	// perform vault operation
	if add {
		rec, err := input(verbose)
		if err != nil {
			log.Fatal(err)
		}
		newRecords := update(rec, records, verbose)
		write(vault, salt, cipher, newRecords, verbose)
		return
	}
	find(vault, salt, cipher, pat, verbose)
}
