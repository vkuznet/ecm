package vault

import (
	"testing"
)

// TestVaultAddRecord function
func TestVaultAddRecord(t *testing.T) {
	vault := Vault{Directory: "Test"}
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
	vault := Vault{Directory: "Test"}
	// perform benchmark test
	for n := 0; n < b.N; n++ {
		_, err := vault.AddRecord("record")
		if err != nil {
			b.Error(err)
		}
	}
}

// TestVaultDeleteRecord function
func TestVaultDeleteRecord(t *testing.T) {
	var records []VaultRecord
	rid := "123"
	rec := VaultRecord{ID: rid}
	records = append(records, rec)
	rid = "567"
	rec = VaultRecord{ID: rid}
	records = append(records, rec)
	vault := Vault{Directory: "Test", Records: records}
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
	var records []VaultRecord
	rid := "123"
	rec := VaultRecord{ID: rid}
	records = append(records, rec)
	// perform benchmark test
	for n := 0; n < b.N; n++ {
		rid = "567"
		rec = VaultRecord{ID: rid}
		records = append(records, rec)
		vault := Vault{Directory: "Test", Records: records}
		err := vault.DeleteRecord(rid)
		if err != nil {
			b.Error(err.Error())
		}
	}
}
