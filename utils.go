package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

const (
	separator = "---\n" // used in pwm data format
)

// StringList implement sort for []string type
type StringList []string

// Len provides length of the []int type
func (s StringList) Len() int { return len(s) }

// Swap implements swap function for []int type
func (s StringList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements less function for []int type
func (s StringList) Less(i, j int) bool { return s[i] < s[j] }

// helper function to determine home area for PWM
func pwmHome() string {
	var err error
	hdir := os.Getenv("PWM_HOME")
	if hdir == "" {
		hdir, err = os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		hdir = fmt.Sprintf("%s/.pwm", hdir)
	}
	return hdir
}

// custom split function based on separator delimiter
func pwmSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

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
	return "\n[green]for help prese [red]Ctrl-H[white]"
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
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	password = strings.Replace(password, "\n", "", -1)
	return password, nil
}
