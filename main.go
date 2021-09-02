package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	uuid "github.com/google/uuid"
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
	var vname string
	flag.StringVar(&vname, "vault", "", "vault name")
	var add bool
	flag.BoolVar(&add, "add", false, "add new record")
	var pat string
	flag.StringVar(&pat, "find", "", "find record pattern")
	var cipher string
	flag.StringVar(&cipher, "cipher", "aes", "cipher to use to initialize the vault")
	var decryptFile string
	flag.StringVar(&decryptFile, "decrypt", "", "decrypt given file name")
	var encryptFile string
	flag.StringVar(&encryptFile, "encrypt", "", "encrypt given file and place it into vault")
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
	if decryptFile != "" {
		password, err := readPassword()
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadFile(decryptFile)
		if err != nil {
			panic(err)
		}
		data, err = decrypt(data, password, cipher)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
		os.Exit(0)
	}

	// parse input config
	configFile := fmt.Sprintf("%s/config.json", pwmHome())
	err := ParseConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// set Theme for our app
	setTheme("grey")

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

	// initialize our vault
	vault := Vault{Cipher: cipher, Secret: salt, Verbose: verbose}
	// create our vault
	err = vault.Create(vname)
	if err != nil {
		fmt.Printf("unable to create vault, error %v", err)
		os.Exit(1)
	}

	// setup logger
	log.SetOutput(new(LogWriter))
	if Config.LogFile != "" {
		logFile := Config.LogFile + "-%Y%m%d"
		rl, err := rotatelogs.New(logFile)
		if err == nil {
			rotlogs := RotateLogWriter{RotateLogs: rl}
			log.SetOutput(rotlogs)
		}
	}

	// encrypt given record
	if encryptFile != "" {
		edata, err := ioutil.ReadFile(encryptFile)
		if err != nil {
			panic(err)
		}
		uid := uuid.NewString()
		attachments := []string{encryptFile}
		rmap := make(Record)
		rec := VaultRecord{ID: uid, Map: rmap, Attachments: attachments}
		data, err = encrypt(edata, vault.Secret, vault.Cipher)
		if err != nil {
			panic(err)
		}
		rec.WriteRecord(vault.Directory, vault.Secret, vault.Cipher, vault.Verbose)
		log.Println("Create new vault record %s", rec.ID)
	}

	// read from our vault
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
	gridView(app, &vault)
}
