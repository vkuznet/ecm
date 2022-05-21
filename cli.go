package main

import (
	"log"
	"os"

	"github.com/atotto/clipboard"
	vt "github.com/vkuznet/ecm/vault"
)

// decrypt record
func decryptFile(dfile, cipher, pcopy string) {
	password, err := readPassword()
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

func cli(vault *vt.Vault,
	cipher, efile, dfile, pat, rid, export, vimport, master, pcopy string, verbose int) {

	// decrypt file if given
	if dfile != "" {
		decryptFile(dfile, cipher, pcopy)
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

	// export vault records
	if export != "" {
		err = vault.Export(export)
		if err != nil {
			log.Fatalf("unable to export vault records, error %v", err)
		}
		os.Exit(0)
	}

	// import records to the vault
	if vimport != "" {
		err = vault.Import(vimport)
		if err != nil {
			log.Fatalf("unable to import records to the vault, error %v", err)
		}
		os.Exit(0)
	}

	// change master password of the vault and re-encrypt all records
	if master != "" {
		err = vault.ChangeMaster(master)
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
