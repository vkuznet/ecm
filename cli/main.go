package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	crypt "github.com/vkuznet/ecm/crypt"
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
	flag.StringVar(&vname, "vault", "", "vault directory name")
	var cipher string
	flag.StringVar(&cipher, "cipher", "", "cipher to use (aes, nacl)")
	var dfile string
	flag.StringVar(&dfile, "decrypt", "", "decrypt given file to stdout")
	var efile string
	flag.StringVar(&efile, "encrypt", "", "encrypt given file and place it into vault")
	var pcopy string
	flag.StringVar(&pcopy, "pcopy", "", "extract given attribute from the record and copy to clipboard")
	var export string
	flag.StringVar(&export, "export", "", "export vault records to given file (ECM JSON native format)")
	var vimport string
	flag.StringVar(&vimport, "import", "", "import records from a given file. Support: CSV, JSON, or ecm.json (native format)")
	var recreate bool
	flag.BoolVar(&recreate, "recreate", false, "recreate vault and its records with new password/cipher")
	var pat string
	flag.StringVar(&pat, "pat", "", "search pattern in vault records")
	var info bool
	flag.BoolVar(&info, "info", false, "show vault info")
	var version bool
	flag.BoolVar(&version, "version", false, "show version")
	var edit string
	flag.StringVar(&edit, "edit", "", "edit record with given ID")
	var add string
	flag.StringVar(&add, "add", "", "add new record (login|card|note|json|file)")
	var rid string
	flag.StringVar(&rid, "rid", "", "show record with given ID and copy its password to clipboard")
	var gen string
	flag.StringVar(&gen, "gen", "", "generate password with given length:attributes (attributes can be 'numbers', symbols' or their combinations), e.g. 16:numbers+symbols will provide password of length 16 with numbers and symbols in it")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbose level")
	var examples bool
	flag.BoolVar(&examples, "examples", false, "show examples")
	flag.Parse()
	if version {
		fmt.Println(ecmInfo())
		os.Exit(0)

	}

	// generate password if asked
	if gen != "" {
		arr := strings.Split(gen, ":")
		i, e := strconv.Atoi(arr[0])
		if e != nil {
			log.Fatal(e)
		}
		var numbers, symbols bool
		if strings.Contains(gen, "numbers") {
			numbers = true
		}
		if strings.Contains(gen, "symbols") {
			symbols = true
		}
		p := crypt.CreatePassword(i, numbers, symbols)
		fmt.Println("New password: ", p)
		os.Exit(0)
	}

	// use file name in a log
	if verbose > 0 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	// initialize our vault
	vault := vt.Vault{Cipher: crypt.GetCipher(cipher), Verbose: verbose, Start: time.Now()}
	if vname == "" {
		// by default vault is located at $HOME/.ecm
		udir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		vname = filepath.Join(udir, ".ecm/Primary")
	}
	vault.Directory = vname
	// create vault if necessary and read its records
	err := vault.Create(vname)
	if err != nil {
		log.Fatalf("unable to create vault, error %v", err)
	}

	// CLI or UI mode
	if examples {
		ecmExamples()
		os.Exit(0)
	}
	cli(
		&vault,
		efile,
		dfile,
		add,
		pat,
		rid,
		edit,
		pcopy,
		export,
		vimport,
		recreate,
		info,
		verbose,
	)
}
