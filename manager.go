package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

func write(vault, secret string) {
	// if vault exists
	//     data, err := read(vault, secret)
	//     if err != nil {
	//         log.Fatal(err)
	//     }

	// replace with input for data record
	vd := VaultData{Name: "test", Value: "value"}
	rec := VaultRecord{Name: "record", Data: []VaultData{vd, vd}, Note: "some note"}

	// our output will consist of set of records
	var records [][]byte

	// marshall single record
	data, err := json.Marshal(rec)
	if err != nil {
		log.Fatal(err)
	}
	// encrypt our record
	data, err = encrypt(data, secret)
	if err != nil {
		log.Fatal(err)
	}
	// add record to the final list of records
	records = append(records, data)

	file, err := os.Create(vault)
	w := bufio.NewWriter(file)
	for _, rec := range records {
		w.Write(rec)
		w.Write([]byte("\n"))
	}
	w.Flush()

	// write our records back to vault
	//     err = ioutil.WriteFile(vault, records, 0777)
	//     if err != nil {
	//         log.Fatalf("Unable to write, file: %s, error: %v\n", vault, err)
	//     }
}

func read(vault, secret string) ([]VaultRecord, error) {
	var records []VaultRecord
	//     data, err := ioutil.ReadFile(vault)
	//     if err != nil {
	//         return records, err
	//     }
	// open file
	f, err := os.Open(vault)
	if err != nil {
		return records, err
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data, err := decrypt([]byte(scanner.Text()), secret)
		if err != nil {
			return records, err
		}
		var rec VaultRecord
		err = json.Unmarshal(data, &rec)
		if err != nil {
			return records, err
		}
		records = append(records, rec)
	}
	return records, nil
}
func find(vault, secret, pat string) {
	records, err := read(vault, secret)
	if err != nil {
		log.Fatal(err)
	}
	for _, rec := range records {
		fmt.Printf("json record %+v", rec)
		printRecord(rec)
	}
}

func printRecord(rec VaultRecord) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Name\t%s\n", rec.Name)
	fmt.Fprintf(w, "Aliases\t%s\n", strings.Join(rec.Aliases, ","))
	fmt.Fprintf(w, "Note\t%s\n", rec.Note)
	fmt.Fprintf(w, "Records:\n")
	for _, r := range rec.Data {
		fmt.Fprintf(w, "%s\t%s\n", r.Name, r.Value)
		fmt.Fprintf(w, "---\n")
	}
	fmt.Fprintf(w, "\n")
}
