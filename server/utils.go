package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vkuznet/ecm/crypt"
	"golang.org/x/crypto/bcrypt"

	// clone of "code.google.com/p/rsc/qr" which no longer available
	// "github.com/vkuznet/rsc/qr"
	qr "rsc.io/qr"

	// imaging library
	"github.com/disintegration/imaging"
)

/*
 * 2fa utils
 */

// helper function to generate QR image file
func generateQRImage(authLink, fname string) error {
	// Encode authLink to QR codes
	// qr.H = 65% redundant level
	// see https://godoc.org/code.google.com/p/rsc/qr#Level
	code, err := qr.Encode(authLink, qr.H)
	if err != nil {
		log.Println("unable to encode auth link", err)
		return err
	}

	imgByte := code.PNG()

	// convert byte to image for saving to file
	img, _, _ := image.Decode(bytes.NewReader(imgByte))

	// TODO: file should be unique for each client
	err = imaging.Save(img, fname)
	if err != nil {
		log.Println("unable to generate QR image file", err)
	}
	return err
}

// getBearerToken returns token from
// HTTP Header "Authorization: Bearer <token>"
func getBearerToken(header string) (string, error) {
	if header == "" {
		return "", fmt.Errorf("An authorization header is required")
	}
	token := strings.Split(header, " ")
	if len(token) != 2 {
		return "", fmt.Errorf("Malformed bearer token")
	}
	return token[1], nil
}

// helper function to check if file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// helper function to generate password hash
func getPasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// helper function to check password hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TokenRecord represents token record
type TokenRecord struct {
	Token  string
	Expire time.Time
}

// helper function to generate and encrypt token. It returns its md5 hash
// and nil error, otherwise it returns empty hash and actual error
func encryptToken(passphrase, cipher string) (string, error) {
	now := time.Now()
	tstamp := now.Format(time.RFC3339Nano)
	token := fmt.Sprintf("token-%s", tstamp)
	expire := now.Add(time.Duration(Config.TokenExpireInterval) * time.Second)
	t := TokenRecord{Token: token, Expire: expire}
	data, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	data, err = crypt.Encrypt(data, passphrase, cipher)
	hash := base64.StdEncoding.EncodeToString(data)
	return hash, err
}

// helper function to decrypt token and check its validity. If eveything is
// fine with our token it returns its hash and nil error, otherwise empty
// hash and actual error is returned
func decryptToken(t, passphrase, cipher string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		return "", err
	}
	data, err = crypt.Decrypt(data, passphrase, cipher)
	if err != nil {
		return "", err
	}
	var trec TokenRecord
	err = json.Unmarshal(data, &trec)
	if err != nil {
		return "", err
	}
	// decode our token
	tstamp := strings.Replace(trec.Token, "token-", "", -1)
	expire := trec.Expire
	ts, err := time.Parse(time.RFC3339Nano, tstamp)
	if err != nil {
		return "", err
	}
	expTime := ts.Add(time.Duration(Config.TokenExpireInterval) * time.Second)
	if expTime != expire {
		return "", errors.New("wrong token expire timestamp")
	}
	return t, err
}
