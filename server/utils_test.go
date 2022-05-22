package main

import (
	"testing"
)

// TestGenToken function
func TestGenToken(t *testing.T) {
	salt := "test"
	cipher := "aes"
	token, err := encryptToken(salt, cipher)
	if err != nil {
		t.Error(err.Error())
	}
	dtoken, err := decryptToken(token, salt, cipher)
	if err != nil {
		t.Error(err.Error())
	}
	if token != dtoken {
		t.Errorf("wrong token, expect '%s' got '%s'", token, dtoken)
	}
}
