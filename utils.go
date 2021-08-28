package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	separator = "---\n" // used in pwm data format
)

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
func backup(src string) (string, int64, error) {
	tstamp := time.Now().Format(time.RFC3339Nano)
	dst := fmt.Sprintf("%s.backup-%s", src, tstamp)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return dst, 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return dst, 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return dst, 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return dst, 0, err
	}
	defer destination.Close()

	err = os.Chmod(dst, 0600)
	if err != nil {
		log.Println("unable to change file permission of", dst)
	}

	nBytes, err := io.Copy(destination, source)
	return dst, nBytes, err
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
	xlist := []string{"", "KB", "MB", "GB", "TB", "PB"}
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
	return "\nfor help prese Ctrl-H"
}

// helper function to return common keys
func helpKeys() string {
	info := "\nCommon keys:"
	info = fmt.Sprintf("%s, Ctrl-F switch to Search", info)
	info = fmt.Sprintf("%s, Ctrl-L switch to Records", info)
	info = fmt.Sprintf("%s, Ctrl-E switch to record edit mode", info)
	info = fmt.Sprintf("%s, Ctrl-N next widget", info)
	info = fmt.Sprintf("%s, Ctrl-P previous widget", info)
	info = fmt.Sprintf("%s, Ctrl-G generate password", info)
	info = fmt.Sprintf("%s, Ctrl-C copy password to clipboard", info)
	info = fmt.Sprintf("%s, Ctrl-Q quit the app", info)
	return info
}
