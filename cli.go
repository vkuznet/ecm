package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	"github.com/vkuznet/ecm/crypt"
	vt "github.com/vkuznet/ecm/vault"
	"golang.org/x/term"
)

// helper function to read vault secret from stdin
func secretPlain(verbose int) (string, error) {
	fmt.Print("\nEnter vault secret: ")
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Println("unable to read stdin, error ", err)
		return "", err
	}
	salt := strings.Replace(string(bytes), "\n", "", -1)
	fmt.Println()
	if verbose > 5 {
		log.Printf("vault secret '%s'", salt)
	}
	return salt, nil
}

// decrypt record
func decryptFile(dfile, cipher, pcopy string) {
	password, err := vt.ReadPassword()
	if err != nil {
		panic(err)
	}
	write := "stdout"
	if pcopy != "" {
		write = "clipboard"
	}
	decryptInput(dfile, password, cipher, write, pcopy)
	os.Exit(0)
}

// cli main function
//gocyclo:ignore
func cli(
	vault *vt.Vault,
	efile, dfile, pat, rid, edit, pcopy, export, vimport string,
	recreate, info bool,
	verbose int,
) {

	// decrypt file if given
	if dfile != "" {
		decryptFile(dfile, vault.Cipher, pcopy)
	}
	// get vault secret
	salt, err := secretPlain(verbose)
	if err != nil {
		log.Fatal(err)
	}
	vault.Secret = salt

	// encrypt given record
	if efile != "" {
		vault.EncryptFile(efile)
		os.Exit(0)
	}

	// read from our vault
	err = vault.Read()
	if err != nil {
		log.Fatal("unable to read vault, error ", err)
	}

	// show vault info
	if info {
		fmt.Println(vault.Info())
		os.Exit(0)
	}

	// edit given record
	if edit != "" {
		err := vault.EditRecord(edit)
		if err != nil {
			log.Fatalf("unable to edit vault record, error '%s'", err)
		}
		os.Exit(0)
	}
	// export vault records
	if export != "" && vimport == "" {
		err = vault.Export(export)
		if err != nil {
			log.Fatalf("unable to export vault records, error %v", err)
		}
		os.Exit(0)
	}

	// import records to the vault
	if vimport != "" {
		err = vault.Import(vimport, export)
		if err != nil {
			log.Fatalf("unable to import records to the vault, error %v", err)
		}
		os.Exit(0)
	}

	// change master password of the vault and re-encrypt all records
	if recreate {
		log.Printf("Supported ciphers: %v", crypt.SupportedCiphers)
		newCipher, err := vt.ReadInput("Cipher to use:")
		if err != nil {
			log.Fatal(err)
		}
		if !InList(newCipher, crypt.SupportedCiphers) {
			log.Fatal("Unsupported cipher")
		}
		newPassword, err := vt.ReadPassword()
		if err != nil {
			log.Fatal(err)
		}
		newPassword2, err := vt.ReadPassword()
		if err != nil {
			log.Fatal(err)
		}
		if newPassword != newPassword2 {
			log.Fatal("provided password strings do not match")
		}
		err = vault.Recreate(newPassword, newCipher)
		if err != nil {
			log.Fatalf("unable to change vault master password and re-encrypt its records, error %v", err)
		}
		os.Exit(0)
	}

	records := vault.Records
	// perform search
	if pat != "" {
		records = vault.Find(pat)
	} else if rid != "" {
		var newRecords []vt.VaultRecord
		// copy record password to clipboard if necessary
		// find give record ID
		for _, rec := range records {
			if rec.ID == rid {
				if pcopy == "" {
					pcopy = "Password" // by default we copy Password to clipboard
				}
				if v, ok := rec.Map[pcopy]; ok {
					if err := clipboard.WriteAll(v); err != nil {
						log.Printf("ERROR: unable to copy '%s' to clipboard", pcopy)
					}
				}
				newRecords = append(newRecords, rec)
				break
			}
		}
		records = newRecords
	}

	// print records
	vt.TabularPrint(records)

}
