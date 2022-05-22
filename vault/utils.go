package vault

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/vkuznet/ecm/crypt"
)

// StringList implement sort for []string type
type StringList []string

// Len provides length of the []int type
func (s StringList) Len() int { return len(s) }

// Swap implements swap function for []int type
func (s StringList) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less implements less function for []int type
func (s StringList) Less(i, j int) bool { return s[i] < s[j] }

// backup helper function to make a vault backup
// based on https://github.com/mactsouk/opensource.com/blob/master/cp1.go
func backup(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		//         log.Printf("file '%s' does not exist, error %v", src, err)
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

// TabularPrint provide tabular print of reocrds
// based on http://networkbit.ch/golang-column-print/
func TabularPrint(records []VaultRecord) {
	// initialize tabwriter
	w := new(tabwriter.Writer)
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()
	for _, rec := range records {
		fmt.Fprintf(w, "\n------------")
		fmt.Fprintf(w, "\nID:\t%s", rec.ID)
		for _, key := range OrderedKeys {
			if val, ok := rec.Map[key]; ok {
				if strings.ToLower(key) == "password" {
					newVal := "*"
					for i := 0; i < len(val); i++ {
						newVal += "*"
					}
					val = newVal
				}
				fmt.Fprintf(w, "\n%v:\t%v", key, val)
			}
		}
		for key, val := range rec.Map {
			if crypt.InList(key, OrderedKeys) {
				continue
			}
			if strings.ToLower(key) == "password" {
				newVal := "*"
				for i := 0; i < len(val); i++ {
					newVal += "*"
				}
				val = newVal
			}
			fmt.Fprintf(w, "\n%v:\t%v", key, val)
		}
		fmt.Fprintf(w, "\n")
	}

}
