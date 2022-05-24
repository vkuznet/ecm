package vault

import (
	"log"
	"os"
	"testing"
	"time"
)

func tempDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	vdir, err := os.MkdirTemp(cwd, "test")
	if err != nil {
		log.Fatal(err)
	}
	return vdir
}

// TestVaultAddRecord function
func TestVaultAddRecord(t *testing.T) {
	vdir := tempDir()
	defer os.RemoveAll(vdir)

	vault := Vault{Directory: vdir, Cipher: "aes", Start: time.Now()}
	vault.AddRecord("record")
	if len(vault.Records) != 1 {
		t.Errorf("fail to add record to the vault")
	}
	rec := vault.Records[0]
	var foundName, foundLogin bool
	for key := range rec.Map {
		if key == "Name" {
			foundName = true
		}
		if key == "Login" {
			foundLogin = true
		}
	}
	if !foundName || !foundLogin {
		t.Error("unable to find name || login in added vault record")
	}

}

// BenchmarkAddRecord provides benchmark test of vault AddRecord functionality
func BenchmarkAddRecord(b *testing.B) {
	vdir := tempDir()
	defer os.RemoveAll(vdir)

	// create new vault in temporary directory
	vault := Vault{Directory: vdir, Cipher: "aes", Start: time.Now()}
	err := vault.Create("TestVault")
	if err != nil {
		b.Fatal(err)
	}

	// perform benchmark test
	for n := 0; n < b.N; n++ {
		// add vault login record
		_, err := vault.AddRecord("login")
		if err != nil {
			b.Error(err)
		}
	}
}

// TestVaultDeleteRecord function
func TestVaultDeleteRecord(t *testing.T) {
	vdir := tempDir()
	defer os.RemoveAll(vdir)

	var records []VaultRecord
	rid := "123"
	rec := VaultRecord{ID: rid}
	records = append(records, rec)
	rid = "567"
	rec = VaultRecord{ID: rid}
	records = append(records, rec)

	vault := Vault{Directory: vdir, Records: records, Cipher: "aes", Start: time.Now()}

	err := vault.DeleteRecord(rid)
	if err != nil {
		t.Error(err.Error())
	}
	found := false
	for _, rec := range records {
		if rec.ID == rid {
			found = true
		}
	}
	if !found {
		t.Errorf("did not find record in a vault")
	}
}

// BenchmarkDeleteRecord provides benchmark test of vault DeleteRecord functionality
func BenchmarkDeleteRecord(b *testing.B) {
	vdir := tempDir()
	defer os.RemoveAll(vdir)

	var records []VaultRecord
	rid := "123"
	rec := VaultRecord{ID: rid}
	records = append(records, rec)

	vault := Vault{Directory: vdir, Records: records, Cipher: "aes", Start: time.Now()}

	// perform benchmark test
	for n := 0; n < b.N; n++ {
		rid = "567"
		rec = VaultRecord{ID: rid}
		records = append(records, rec)
		vault.Records = records
		err := vault.DeleteRecord(rid)
		if err != nil {
			b.Error(err.Error())
		}
	}
}
