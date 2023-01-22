package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/crypto/ssh/terminal"
)

// main function
func main() {
	usr, _ := user.Current()
	defaultPath := fmt.Sprintf("%s/.keepass.kdbx", usr.HomeDir)

	var kfile, kpath string
	flag.StringVar(&kpath, "kdbx", defaultPath, "path to kdbx file")
	flag.StringVar(&kfile, "kfile", "", "key file name")
	var interval int
	flag.IntVar(&interval, "interval", 30, "timeout interval in seconds")
	flag.Parse()

	file, err := os.Open(kpath)
	if err != nil {
		log.Fatal(err)
	}

	db := gokeepasslib.NewDatabase()
	pwd := getpass("Database Password: ")
	if kfile != "" {
		db.Credentials, err = gokeepasslib.NewPasswordAndKeyCredentials(pwd, kfile)
		if err != nil {
			log.Fatalf("ERROR: unable to get credentials, %v", err)
		}
	} else {
		db.Credentials = gokeepasslib.NewPasswordCredentials(pwd)
	}
	_ = gokeepasslib.NewDecoder(file).Decode(db)
	db.UnlockProtectedEntries()

	time0 := time.Now()
	timeout := time.Duration(interval) * time.Second

	searchMsg := "\nsearch for: "
	fmt.Printf(searchMsg)

	// we'll read out std input via goroutine
	ch := make(chan string)
	go readInputChannel(ch)

	// read stdin and search for DB record
	for {
		select {
		case input := <-ch:
			input = strings.Replace(input, "\n", "", -1)
			search(db, input)
			time0 = time.Now()
			fmt.Printf(searchMsg)
		default:
			if time.Since(time0) > timeout {
				fmt.Printf("\nExit after %s of inactivity", time.Since(time0))
				os.Exit(1)
			}
			time.Sleep(time.Duration(100) * time.Millisecond) // wait for new input
		}
	}
}

// helper function to get value of kdbx record
func getValue(entry gokeepasslib.Entry, key string) string {
	if ptr := entry.Get(key); ptr != nil {
		if ptr.Key == key {
			return fmt.Sprintf("%+v", ptr.Value.Content)
		}
	}
	return ""
}

// helper function to read stdin and send it over provided channel
// https://stackoverflow.com/questions/50788805/how-to-read-from-stdin-with-goroutines-in-golang
func readInputChannel(ch chan<- string) {
	for {
		val, err := readInput()
		if err != nil {
			log.Println("WARNING: wrong input", err)
		}
		ch <- val
	}
}

// helper function to read from stdin
func readInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	val, err := reader.ReadString('\n')
	return val, err
}

// helper function to search for given input
func search(db *gokeepasslib.Database, input string) {
	rsearch, _ := regexp.Compile("(?i)" + input)
	found := make(map[string]string)

	for _, top := range db.Content.Root.Groups {
		if top.Name == "NewDatabase" {
			fmt.Println("wrong password or empty database")
			os.Exit(3)
		}
		for _, groups := range top.Groups {
			for _, entry := range groups.Entries {
				entry_path := fmt.Sprintf("%s/%s/%s", top.Name, groups.Name, entry.GetTitle())
				if strings.Compare(entry.GetTitle(), input) == 0 {
					printRecord(entry)
					found[entry_path] = entry.GetPassword()
				} else if rsearch.MatchString(entry_path) {
					printRecord(entry)
					found[entry_path] = entry.GetPassword()
				}
			}
		}
	}

	if len(found) == 1 {
		for key, found_pw := range found {
			if err := clipboard.WriteAll(string(found_pw)); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s password copied to clipboard\n", key)
		}
	}
}

// helper function to print record
func printRecord(entry gokeepasslib.Entry) {
	fmt.Printf("Title    %s\n", getValue(entry, "Title"))
	fmt.Printf("UserName %s\n", getValue(entry, "UserName"))
	fmt.Printf("URL      %s\n", getValue(entry, "URL"))
	fmt.Printf("Notes    %s\n", getValue(entry, "Notes"))
	fmt.Printf("Tags     %s\n", entry.Tags)
}

// helper function to get password from stdin
func getpass(msg string) string {
	fmt.Print(msg)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err == nil {
		fmt.Println("")
	} else {
		fmt.Println("\nError in ReadPassword", err)
		os.Exit(1)
	}
	password := string(bytePassword)

	return strings.TrimSpace(password)
}
