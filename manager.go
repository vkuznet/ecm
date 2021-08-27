package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/google/uuid"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

// helper function to read vault secret from stdin
func secret(verbose int) (string, error) {
	fmt.Print("Enter vault secret: ")
	bytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Println("unable to read stdin, error ", err)
		return "", err
	}
	salt := strings.Replace(string(bytes), "\n", "", -1)
	fmt.Println()
	if verbose > 0 {
		log.Printf("vault secret '%s'", salt)
	}
	return salt, nil
}

// helper function to get user input
func input(verbose int) (VaultRecord, error) {
	app := tview.NewApplication()
	rec := inputForm(app)
	return rec, nil
}

// helper function to get user input
func inputOld(verbose int) (VaultRecord, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nEnter record name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return VaultRecord{}, err
	}
	name = strings.Replace(name, "\n", " ", -1)

	fmt.Print("\nEnter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return VaultRecord{}, err
	}
	username = strings.Replace(username, "\n", "", -1)

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return VaultRecord{}, err
	}
	password := string(bytePassword)
	password = strings.Replace(password, "\n", "", -1)

	fmt.Print("\nEnter URL: ")
	rurl, err := reader.ReadString('\n')
	if err != nil {
		return VaultRecord{}, err
	}
	rurl = strings.Replace(rurl, "\n", " ", -1)

	fmt.Print("\nEnter note: ")
	note, err := reader.ReadString('\n')
	if err != nil {
		return VaultRecord{}, err
	}
	note = strings.Replace(note, "\n", " ", -1)

	// replace with input for data record
	recLogin := VaultItem{Name: "login", Value: username}
	recPassword := VaultItem{Name: "password", Value: password}
	recUrl := VaultItem{Name: "url", Value: rurl}
	uid := uuid.NewString()
	rec := VaultRecord{ID: uid, Name: name, Items: []VaultItem{recLogin, recPassword, recUrl}, Note: note}
	return rec, nil
}

// helper function to read vault and return list of records
func read(vault, secret, cipher string, verbose int) ([]VaultRecord, error) {
	var records []VaultRecord

	// check first if file exsist
	if _, err := os.Stat(vault); os.IsNotExist(err) {
		return records, nil
	}

	// open file
	file, err := os.Open(vault)
	if err != nil {
		log.Println("unable to open a vault", err)
		return records, err
	}
	// remember to close the file at the end of the program
	defer file.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(file)
	scanner.Split(pwmSplitFunc)
	for scanner.Scan() {
		text := scanner.Text()
		textData := []byte(text)
		if verbose > 0 {
			log.Printf("read record\n%v\n", textData)
		}

		data := textData
		if cipher != "" {
			data, err = decrypt(textData, secret, cipher)
			if err != nil {
				log.Printf("unable to decrypt data\n%v\nerror %v", textData, err)
				return records, err
			}
		}

		var rec VaultRecord
		err = json.Unmarshal(data, &rec)
		if err != nil {
			log.Println("ERROR: unable to unmarshal the data", err)
			return records, err
		}
		records = append(records, rec)
	}
	return records, nil
}

// helper function to find given pattern in vault records
func find(vault, secret, cipher, pat string, verbose int) {
	records, err := read(vault, secret, cipher, verbose)
	if err != nil {
		log.Fatal(err)
	}
	for _, rec := range records {
		if verbose > 0 {
			fmt.Printf("json record %+v", rec)
		}
		printRecord(rec)
	}
}

// helper function to print the record
func printRecord(rec VaultRecord) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Name\t%s\n", rec.Name)
	fmt.Fprintf(w, "URL\t%s\n", rec.URL)
	fmt.Fprintf(w, "Tags\t%s\n", strings.Join(rec.Tags, ","))
	fmt.Fprintf(w, "Note\t%s\n", rec.Note)
	fmt.Fprintf(w, "Records:\n")
	for _, r := range rec.Items {
		fmt.Fprintf(w, "%s\t\t%s\n", r.Name, r.Value)
		fmt.Fprintf(w, "---\n")
	}
	fmt.Fprintf(w, "\n")
}
