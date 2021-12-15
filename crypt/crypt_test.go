package crypt

import (
	"testing"

	"github.com/vkuznet/gpm/crypt"
)

// CryptEncodeDecode function
func CryptEncodeDecode(t *testing.T) {
	salt := "test"
	data := []byte(salt)
	ciphers := []string{"aes", "nacl"}
	for _, c := range ciphers {
		edata, err := crypt.Encrypt(data, salt, c)
		if err != nil {
			t.Error(err.Error())
		}
		result, err := crypt.Decrypt(edata, salt, c)
		if string(edata) != string(result) {
			t.Errorf("encrypt/decrypt failure with %s cipher", c)
		}
		if err != nil {
			t.Error(err.Error())
		}
	}
	_, err := crypt.Encrypt(data, salt, "lsdjflksdj")
	if err == nil {
		t.Error("failt to recognize unsupported cipher")
	}
}

// BenchmarkEncryptAES provides benchmark test for AES encrypt operation
func BenchmarkEncryptAES(b *testing.B) {
	salt := "test"
	data := []byte(salt)
	cipher := "aes"
	for n := 0; n < b.N; n++ {
		_, err := crypt.Encrypt(data, salt, cipher)
		if err != nil {
			b.Error(err.Error())
		}
	}
}

// BenchmarkEncryptNaCl provides benchmark test for NaCl encrypt operation
func BenchmarkEncryptNaCl(b *testing.B) {
	salt := "test"
	data := []byte(salt)
	cipher := "nacl"
	for n := 0; n < b.N; n++ {
		_, err := crypt.Encrypt(data, salt, cipher)
		if err != nil {
			b.Error(err.Error())
		}
	}
}

// BenchmarkDecryptAES provides benchmark test for AES encrypt operation
func BenchmarkDecryptAES(b *testing.B) {
	salt := "test"
	data := []byte(salt)
	cipher := "aes"
	edata, err := crypt.Encrypt(data, salt, cipher)
	if err != nil {
		b.Error(err.Error())
	}
	for n := 0; n < b.N; n++ {
		res, err := crypt.Decrypt(edata, salt, cipher)
		if err != nil {
			b.Error(err.Error())
		}
		if string(res) != string(data) {
			b.Error("fail to decrypt data with aes cipher")
		}
	}
}

// BenchmarkDecryptNaCl provides benchmark test for NaCl encrypt operation
func BenchmarkDecryptNaCl(b *testing.B) {
	salt := "test"
	data := []byte(salt)
	cipher := "nacl"
	edata, err := crypt.Encrypt(data, salt, cipher)
	if err != nil {
		b.Error(err.Error())
	}
	for n := 0; n < b.N; n++ {
		res, err := crypt.Decrypt(edata, salt, cipher)
		if err != nil {
			b.Error(err.Error())
		}
		if string(res) != string(data) {
			b.Error("fail to decrypt data with nacl cipher")
		}
	}
}
