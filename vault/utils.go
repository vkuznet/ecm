package vault

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	utils "github.com/vkuznet/ecm/utils"
)

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
			if utils.InList(key, OrderedKeys) {
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

// helper function to return black message on white bold foreground
func saveMessage(msg string) string {
	c := color.New(color.FgBlack).Add(color.BgWhite).Add(color.Bold)
	return c.Sprint(msg)
}
