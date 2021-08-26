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

// update given records in our vault
func update(rec VaultRecord, records []VaultRecord, verbose int) []VaultRecord {
	// TODO: so far we add record to the list of records
	// but I need to implement searhc for that record and if there is one
	// we need to update it in place

	// add record to the final list of records
	records = append(records, rec)
	return records
}

// helper function to write vault record
func write(vault, secret, cipher string, records []VaultRecord, verbose int) {
	// make backup vault first
	backup(vault)

	file, err := os.Create(vault)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)
	for _, rec := range records {
		// marshall single record
		data, err := json.Marshal(rec)
		if err != nil {
			log.Fatal(err)
		}

		// encrypt our record
		if verbose > 0 {
			log.Printf("record '%s' secret '%s'\n", string(data), secret)
		}
		edata := data
		if cipher != "" {
			edata, err = encrypt(data, secret, cipher)
			if err != nil {
				log.Fatal(err)
			}
		}
		if verbose > 1 {
			log.Printf("write data record\n%v\nsecret %v", edata, secret)
		}
		w.Write(edata)
		w.Write([]byte("---\n"))
		w.Flush()
	}

	// write our records back to vault
	//     err = ioutil.WriteFile(vault, records, 0777)
	//     if err != nil {
	//         log.Fatalf("Unable to write, file: %s, error: %v\n", vault, err)
	//     }
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
		//         if verbose > 0 {
		//             log.Printf("read '%v'", text)
		//         }
		textData := []byte(text)
		if verbose > 0 {
			log.Printf("read record\n%v\nsecret %v", textData, secret)
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
