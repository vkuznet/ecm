package main

// https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"strings"

	"golang.org/x/crypto/nacl/secretbox"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encryptAES(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
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

func decryptAES(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
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

// the follow code is borrowed from gocrypt repo
// https://github.com/kisom/gocrypto
// see Chapter 3

const (
	// KeySize is the size of a NaCl secret key.
	KeySize = 32

	// NonceSize is the size of a NaCl nonce.
	NonceSize = 24
)

// GenerateKey creates a new secret key either randomly if input key is
// not provided or via key hash
func GenerateKey(passphrase string) (*[KeySize]byte, error) {
	key := new([KeySize]byte)
	if passphrase != "" {
		hash := []byte(createHash(passphrase))
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

// Encrypt generates a random nonce and encrypts the input using
// NaCl's secretbox package. The nonce is prepended to the ciphertext.
// A sealed message will the same size as the original message plus
// secretbox.Overhead bytes long.
func Encrypt(key *[KeySize]byte, message []byte) ([]byte, error) {
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, ErrEncrypt
	}

	out := make([]byte, len(nonce))
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce, key)
	return out, nil
}

// Decrypt extracts the nonce from the ciphertext, and attempts to
// decrypt with NaCl's secretbox.
func Decrypt(key *[KeySize]byte, message []byte) ([]byte, error) {
	if len(message) < (NonceSize + secretbox.Overhead) {
		log.Println("message length is less than nonce size+overhead")
		return nil, ErrDecrypt
	}

	var nonce [NonceSize]byte
	copy(nonce[:], message[:NonceSize])
	out, ok := secretbox.Open(nil, message[NonceSize:], &nonce, key)
	if !ok {
		log.Println("fail to open secret box")
		return nil, ErrDecrypt
	}

	return out, nil
}

// our encrypt wrapper function
func encrypt(data []byte, passphrase, cipher string) ([]byte, error) {
	if strings.ToLower(cipher) == "naci" {
		key, err := GenerateKey(passphrase)
		if err != nil {
			return []byte{}, err
		}
		return Encrypt(key, data)
	}
	return encryptAES(data, passphrase)
}

// our decrypt wrapper function
func decrypt(data []byte, passphrase, cipher string) ([]byte, error) {
	if strings.ToLower(cipher) == "naci" {
		key, err := GenerateKey(passphrase)
		if err != nil {
			return []byte{}, err
		}
		return Decrypt(key, data)
	}
	return decryptAES(data, passphrase)
}
