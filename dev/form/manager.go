package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rivo/tview"
	vt "github.com/vkuznet/ecm/vault"
)

// helper function to read vault secret from stdin
func secret(verbose int) (string, error) {
	app := tview.NewApplication()
	salt, err := lockView(app, verbose)
	app.Stop()
	return salt, err
}

// helper function to get user input
func input(verbose int) (vt.VaultRecord, error) {
	app := tview.NewApplication()
	rec := inputForm(app)
	return rec, nil
}

// helper function to print the record
func printRecord(rec vt.VaultRecord) {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer w.Flush()
	fmt.Fprintf(w, "\n")
	for key, val := range rec.Map {
		fmt.Fprintf(w, "%s\t%s\n", key, val)
	}
	fmt.Fprintf(w, "\n")
}
