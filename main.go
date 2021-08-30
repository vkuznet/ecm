package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
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
	var fname string
	flag.StringVar(&fname, "vault", "", "vault file name")
	var add bool
	flag.BoolVar(&add, "add", false, "add new record")
	var pat string
	flag.StringVar(&pat, "find", "", "find record pattern")
	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbose level")
	flag.Parse()
	if version {
		fmt.Println(info())
		os.Exit(0)

	}

	// parse input config
	configFile := fmt.Sprintf("%s/config.json", pwmHome())
	err := ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// log time, filename, and line number
	if verbose > 0 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// setup logger
	log.SetOutput(new(LogWriter))
	if Config.LogFile != "" {
		rl, err := rotatelogs.New(Config.LogFile + "-%Y%m%d")
		if err == nil {
			rotlogs := RotateLogWriter{RotateLogs: rl}
			log.SetOutput(rotlogs)
		}
	}

	// check if fname is relative and construct proper full path
	if fname != "" {
		abs, err := filepath.Abs(fname)
		if err != nil {
			log.Fatal(err)
		}
		fname = abs
	}

	// determine vault location and if it is not provided or does not exists
	// creat $HOME/.pwm area and assign new vault file there
	_, err = os.Stat(fname)
	if fname == "" || os.IsNotExist(err) {
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
		fname = fmt.Sprintf("%s/vault.%s", vdir, Config.Cipher)
	}

	// get vault secret
	salt, err := secret(verbose)
	if err != nil {
		log.Fatal(err)
	}

	// initialize our vault
	vault := Vault{Filename: fname, Cipher: Config.Cipher, Secret: salt, Verbose: verbose}
	if verbose > 0 {
		log.Println(vault.Info())
	}
	err = vault.Read()
	if err != nil {
		log.Fatal("unable to read vault, error ", err)
	}

	// perform vault operation
	if add {
		rec, err := input(verbose)
		if err != nil {
			log.Fatal(err)
		}
		if verbose > 0 {
			log.Println("get new input record", rec.String())
		}
		vault.Update(rec)
		vault.Write()
	}

	app := tview.NewApplication()
	//     listForm(app, records)
	gridView(app, &vault)

	//     find(vault, salt, cipher, pat, verbose)
}
