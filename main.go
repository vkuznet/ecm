package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	//     "github.com/rivo/tview"
)

// version of the code
var gitVersion, gitTag string

// Info function returns version string of the server
func info() string {
	goVersion := runtime.Version()
	tstamp := time.Now().Format("2006-02-01")
	return fmt.Sprintf("gpm git=%s tag=%s go=%s date=%s", gitVersion, gitTag, goVersion, tstamp)
}

func main() {
	var vname string
	flag.StringVar(&vname, "vault", "", "vault name")
	var cipher string
	flag.StringVar(&cipher, "cipher", "", "cipher to use (aes, nacl)")
	var dfile string
	flag.StringVar(&dfile, "decrypt", "", "decrypt given file name")
	var efile string
	flag.StringVar(&efile, "encrypt", "", "encrypt given file and place it into vault")
	var attr string
	flag.StringVar(&attr, "attr", "", "extract certain attribute from the record")
	var write string
	flag.StringVar(&write, "write", "stdout", "write record to (stdout|clipboard|<filename>)")
	var export string
	flag.StringVar(&export, "export", "", "export vault records to given file")
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var lockInterval int
	flag.IntVar(&lockInterval, "lock", 60, "lock interval in seconds")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbose level")
	var serverConfig string
	flag.StringVar(&serverConfig, "serverConfig", "", "start HTTP server with provided configuration")
	flag.Parse()
	if version {
		fmt.Println(info())
		os.Exit(0)

	}

	// use file name in a log
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// start HTTP server if it is required
	if serverConfig != "" {
		startServer(serverConfig)
		os.Exit(0)
	}

	// decrypt record
	if dfile != "" {
		password, err := readPassword()
		if err != nil {
			panic(err)
		}
		decryptInput(dfile, password, cipher, write, attr)
		os.Exit(0)
	}

	// parse input config
	configFile := fmt.Sprintf("%s/config.json", gpmHome())
	err := ParseConfig(configFile, verbose)
	if err != nil {
		log.Fatal(err)
	}

	// initialize our vault
	vault := Vault{Cipher: getCipher(cipher), Verbose: verbose, Start: time.Now()}

	// create our vault
	err = vault.Create(vname)
	if err != nil {
		log.Fatalf("unable to create vault, error %v", err)
	}

	// we split either at CLI or UI mode
	if efile != "" || export != "" {
		// get vault secret
		salt, err := secret(verbose)
		if err != nil {
			log.Fatal(err)
		}
		vault.Secret = salt

		// encrypt given record
		if efile != "" {
			vault.EncryptFile(efile)
		}

		// read from our vault
		err = vault.Read()
		if err != nil {
			log.Fatal("unable to read vault, error ", err)
		}

		// export vault records
		if export != "" {
			err = vault.Export(export)
			if err != nil {
				log.Fatalf("unable to export vault records, error %v", err)
			}
			os.Exit(0)
		}
	} else { // UI mode
		setTheme("grey")
		gpgApp(&vault, lockInterval)
	}
}
