package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/vkuznet/ecm/crypt"
	utils "github.com/vkuznet/ecm/utils"
	vt "github.com/vkuznet/ecm/vault"
)

const (
	separator = "---\n" // used in ecm data format
)

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
		log.Printf("file '%s' does not exist, error %v", src, err)
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

// helper function to extract cipher name from file extension
func fileCipher(fname string) string {
	arr := strings.Split(fname, ".")
	cipher := strings.Split(arr[len(arr)-1], "-")[0]
	if !utils.InList(cipher, crypt.SupportedCiphers) {
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
	if !utils.InList(cipher, crypt.SupportedCiphers) {
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
		var rec vt.VaultRecord
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
