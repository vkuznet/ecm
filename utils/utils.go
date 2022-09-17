package utils

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/errors"
	"golang.org/x/term"
)

// StringList implement sort for []string type
type StringList []string

// Len provides length of the []int type
func (s StringList) Len() int { return len(s) }

// Swap implements swap function for []int type
func (s StringList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements less function for []int type
func (s StringList) Less(i, j int) bool { return s[i] < s[j] }

// EcmHome returns home area for ECM
func EcmHome() string {
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

// InList checks item in a list
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

// RandomString generates random string
func RandomString(strSize int, randType string) string {
	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "alpha" {
		dictionary = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	if randType == "number" {
		dictionary = "0123456789"
	}

	var bytes = make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

// MacAddress returns mac address of the host
// https://gist.github.com/tsilvers/085c5f39430ced605d970094edf167ba
func MacAddress() uint64 {
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

// FileExists checks if file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// GetPasswordHash generates password hash
func GetPasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks password hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ReadInput read user input from stdout
// https://gosamples.dev/read-user-input/
func ReadInput(msg string) (string, error) {
	fmt.Print(msg)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		return "", err
	}
	return scanner.Text(), nil
}

// ReadPassword reads input password from stdout
func ReadPassword() (string, error) {
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
