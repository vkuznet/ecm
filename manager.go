package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"text/tabwriter"

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
	if verbose > 5 {
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
