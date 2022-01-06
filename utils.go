package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	"github.com/vkuznet/ecm/crypt"
	"golang.org/x/term"
)

const (
	separator = "---\n" // used in ecm data format
)

// StringList implement sort for []string type
type StringList []string

// Len provides length of the []int type
func (s StringList) Len() int { return len(s) }

// Swap implements swap function for []int type
func (s StringList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements less function for []int type
func (s StringList) Less(i, j int) bool { return s[i] < s[j] }

// helper function to determine home area for ECM
func ecmHome() string {
	var err error
	hdir := os.Getenv("ECM_HOME")
	if hdir == "" {
		hdir, err = os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		hdir = fmt.Sprintf("%s/.ecm", hdir)
	}
	return hdir
}

// custom split function based on separator delimiter
func ecmSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := strings.Index(string(data), separator); i >= 0 {
		return i + len(separator), data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return
}

// backup helper function to make a vault backup
// based on https://github.com/mactsouk/opensource.com/blob/master/cp1.go
func backup(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		log.Println("file src does not exist, error ", err)
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	err = os.Chmod(dst, 0600)
	if err != nil {
		log.Println("unable to change file permission of", dst)
	}

	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// InList helper function to check item in a list
func InList(a string, list []string) bool {
	check := 0
	for _, b := range list {
		if b == a {
			check += 1
		}
	}
	if check != 0 {
		return true
	}
	return false
}

// SizeFormat helper function to convert size into human readable form
func SizeFormat(val interface{}) string {
	var size float64
	var err error
	switch v := val.(type) {
	case int:
		size = float64(v)
	case int32:
		size = float64(v)
	case int64:
		size = float64(v)
	case float64:
		size = v
	case string:
		size, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
	default:
		return fmt.Sprintf("%v", val)
	}
	base := 1000. // CMS convert is to use power of 10
	xlist := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	for _, vvv := range xlist {
		if size < base {
			return fmt.Sprintf("%v (%3.1f%s)", val, size, vvv)
		}
		size = size / base
	}
	return fmt.Sprintf("%v (%3.1f%s)", val, size, xlist[len(xlist)])
}

// helper function to make message about help key
func helpKey() string {
	return "\n[green]for help press [red]Ctrl-H[white]"
}

// helper function to return common keys
func helpKeys() string {
	info := "\nCommon keys:\n"
	info = fmt.Sprintf("%s, [red]Ctrl-N[white] next widget", info)
	info = fmt.Sprintf("%s, [red]Ctrl-P[white] previous widget", info)
	info = fmt.Sprintf("%s\n", info)
	info = fmt.Sprintf("%s, [red]Ctrl-F[white] switch to Search", info)
	info = fmt.Sprintf("%s, [red]Ctrl-L[white] switch to Records", info)
	info = fmt.Sprintf("%s, [red]Ctrl-E[white] record edit mode", info)
	info = fmt.Sprintf("%s\n", info)
	info = fmt.Sprintf("%s, [red]Ctrl-G[white] generate password", info)
	info = fmt.Sprintf("%s, [red]Ctrl-P[white] copy password to clipboard", info)
	info = fmt.Sprintf("%s\n", info)
	info = fmt.Sprintf("%s, [red]Ctrl-Q[white] Exit", info)
	return info
}

// helper function to read input password
func readPassword() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println("")
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	password = strings.Replace(password, "\n", "", -1)
	return password, nil
}

// helper function to extract cipher name from file extension
func fileCipher(fname string) string {
	arr := strings.Split(fname, ".")
	cipher := strings.Split(arr[len(arr)-1], "-")[0]
	if !InList(cipher, crypt.SupportedCiphers) {
		log.Fatalf("given cipher %s is not supported, please use one from the following %v", cipher, crypt.SupportedCiphers)
	}
	return cipher
}

// helper function to decrypt given input (file or stdin)
func decryptInput(fname, password, cipher, write, attr string) {
	var err error
	if cipher == "" {
		cipher = fileCipher(fname)
	}
	if !InList(cipher, crypt.SupportedCiphers) {
		log.Fatalf("given cipher %s is not supported, please use one from the following %v", cipher, crypt.SupportedCiphers)
	}
	var data []byte
	if fname == "-" { // stdin
		var input string
		reader := bufio.NewReader(os.Stdin)
		input, err = reader.ReadString('\n')
		data = []byte(input)
	} else {
		data, err = os.ReadFile(fname)
	}
	if err != nil {
		panic(err)
	}
	data, err = crypt.Decrypt(data, password, cipher)
	if err != nil {
		log.Fatal(err)
	}
	if attr != "" {
		var rec VaultRecord
		err := json.Unmarshal(data, &rec)
		if err != nil {
			log.Fatal("unable to unarmashal vault record", err)
		}
		val, ok := rec.Map[attr]
		if ok {
			data = []byte(val)
		} else {
			log.Fatalf("unable to extract attribute '%s' from the record %+v", attr, rec)
		}
	}
	if write == "stdout" {
		fmt.Println(string(data))
	} else if write == "clipboard" {
		if err := clipboard.WriteAll(string(data)); err != nil {
			log.Fatal("unable to copy to clipboard, error", err)
		}
	} else {
		err := os.WriteFile(write, data, 0755)
		if err != nil {
			log.Fatal("unable to data to output file", err)
		}
	}
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

// getCipher returns either default or given cipher
func getCipher(cipher string) string {
	if cipher == "" {
		cipher = crypt.SupportedCiphers[0]
	}
	if !InList(cipher, crypt.SupportedCiphers) {
		log.Fatalf("given cipher %s is not supported, please use one from the following %v", cipher, crypt.SupportedCiphers)
	}
	return strings.ToLower(cipher)
}

// https://gist.github.com/tsilvers/085c5f39430ced605d970094edf167ba
func macAddress() uint64 {
	interfaces, err := net.Interfaces()
	if err != nil {
		return uint64(0)
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
			// Skip locally administered addresses
			if i.HardwareAddr[0]&2 == 2 {
				continue
			}
			var mac uint64
			for j, b := range i.HardwareAddr {
				if j >= 8 {
					break
				}
				mac <<= 8
				mac += uint64(b)
			}
			return mac
		}
	}
	return uint64(0)
}
