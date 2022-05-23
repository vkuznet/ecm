package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
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
	flag.StringVar(&vname, "vault", "", "vault name")
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
	var rid string
	flag.StringVar(&rid, "rid", "", "show record with given ID and copy its password to clipboard")
	var lockInterval int
	flag.IntVar(&lockInterval, "lock", 60, "lock interval in seconds")
	var verbose int
	flag.IntVar(&verbose, "verbose", 0, "verbose level")
	//     var serverConfig string
	//     flag.StringVar(&serverConfig, "serverConfig", "", "start HTTP server with provided configuration")
	var examples bool
	flag.BoolVar(&examples, "examples", false, "show examples")
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
	vault := vt.Vault{Cipher: crypt.GetCipher(cipher), Verbose: verbose, Start: time.Now()}

	// create vault if necessary
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
