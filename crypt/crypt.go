package crypt

// crypt module provides various ciphers
// for more information see
// https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/
// https://github.com/kisom/gocrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	utils "github.com/vkuznet/ecm/utils"
	"golang.org/x/crypto/nacl/secretbox"
)

// CreateHash creates a hash for given key
func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// SupportedCiphers provides list of supported ciphers
var SupportedCiphers = []string{"aes", "nacl"}

// Cipher defines cipher interface
type Cipher interface {
	Encript(data []byte, key string) ([]byte, error)
	Decript(data []byte, key string) ([]byte, error)
}

// CipherAES represents AES Cipher
type CipherAES struct {
}

// Encrypt implementation for AES cipher
func (c *CipherAES) Encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt implementation for AES Cipher
func (c *CipherAES) Decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return []byte{}, err
	}
	return plaintext, nil
}

const (
	// KeySize is the size of a NaCl secret key
	KeySize = 32
	// NonceSize is the size of a NaCl nonce
	NonceSize = 24
)

// GenerateKey creates a new secret key either randomly if input key is
// not provided or via key hash
func GenerateKey(passphrase string) (*[KeySize]byte, error) {
	key := new([KeySize]byte)
	if passphrase != "" {
		hash := []byte(CreateHash(passphrase))
		for i, v := range hash {
			if i < KeySize {
				key[i] = v
			}
		}
		return key, nil
	}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}

	return key, nil
}

// GenerateNonce creates a new random nonce.
func GenerateNonce() (*[NonceSize]byte, error) {
	nonce := new([NonceSize]byte)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return nonce, nil
}

var (
	// ErrEncrypt is returned when encryption fails.
	ErrEncrypt = errors.New("secret: encryption failed")

	// ErrDecrypt is returned when decryption fails.
	ErrDecrypt = errors.New("secret: decryption failed")
)

// CipherNaCl represents NaCl Cipher
type CipherNaCl struct {
}

// Encrypt implementation of NaCl cipher
func (c *CipherNaCl) Encrypt(data []byte, passphrase string) ([]byte, error) {
	key, err := GenerateKey(passphrase)
	if err != nil {
		return []byte{}, err
	}
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, ErrEncrypt
	}

	out := make([]byte, len(nonce))
	copy(out, nonce[:])
	out = secretbox.Seal(out, data, nonce, key)
	return out, nil
}

// Decrypt implementation of NaCl cipher
func (c *CipherNaCl) Decrypt(data []byte, passphrase string) ([]byte, error) {
	key, err := GenerateKey(passphrase)
	if err != nil {
		return []byte{}, err
	}
	if len(data) < (NonceSize + secretbox.Overhead) {
		log.Println("message length is less than nonce size+overhead")
		return nil, ErrDecrypt
	}

	var nonce [NonceSize]byte
	copy(nonce[:], data[:NonceSize])
	out, ok := secretbox.Open(nil, data[NonceSize:], &nonce, key)
	if !ok {
		log.Println("fail to open secret box")
		return nil, ErrDecrypt
	}

	return out, nil
}

// Encrypt wrapper function to encrypt given binary data blob using given passphrase and cipher
func Encrypt(data []byte, passphrase, cipher string) ([]byte, error) {
	if strings.ToLower(cipher) == "nacl" {
		c := CipherNaCl{}
		return c.Encrypt(data, passphrase)
	} else if strings.ToLower(cipher) == "aes" {
		c := CipherAES{}
		return c.Encrypt(data, passphrase)
	}
	msg := fmt.Sprintf("unsupported cipher %s", cipher)
	return []byte{}, errors.New(msg)
}

// Decrypt wrapper function to decrypt given binary data blob using given passphrase and cipher
func Decrypt(data []byte, passphrase, cipher string) ([]byte, error) {
	if strings.ToLower(cipher) == "nacl" {
		c := CipherNaCl{}
		return c.Decrypt(data, passphrase)
	} else if strings.ToLower(cipher) == "aes" {
		c := CipherAES{}
		return c.Decrypt(data, passphrase)
	}
	msg := fmt.Sprintf("unsupported cipher %s", cipher)
	return []byte{}, errors.New(msg)
}

// GetCipher returns either default or given cipher
func GetCipher(cipher string) string {
	if cipher == "" {
		cipher = SupportedCiphers[0]
	}
	if !utils.InList(cipher, SupportedCiphers) {
		log.Fatalf(
			"given cipher %s is not supported, please use one from the following %v",
			cipher,
			SupportedCiphers,
		)
	}
	return strings.ToLower(cipher)
}
