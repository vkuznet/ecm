package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/atotto/clipboard"
	gokeepasslib "github.com/tobischo/gokeepasslib/v3"
	"golang.org/x/crypto/ssh/terminal"
)

// DBRecords defines map of DB records
type DBRecords map[int]gokeepasslib.Entry

// main function
func main() {
	usr, _ := user.Current()
	defaultPath := fmt.Sprintf("%s/.keepass.kdbx", usr.HomeDir)

	var kfile, kpath string
	flag.StringVar(&kpath, "kdbx", defaultPath, "path to kdbx file")
	flag.StringVar(&kfile, "kfile", "", "key file name")
	var interval int
	flag.IntVar(&interval, "interval", 30, "timeout interval in seconds")
	flag.Usage = func() {
		fmt.Println("Usage: kpass [options]")
		flag.PrintDefaults()
		fmt.Println("Commands within DB")
		fmt.Println("cp <ID> <attribute>   # to copy record ID attribute to cpilboard")
		fmt.Println("rm <ID>               # to remove record ID from database")
		fmt.Println("add <login|note|card> # to add specific record type")
		fmt.Println("timeout <int>         # set timeout interval in seconds")
	}
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

	searchMsg := "\ncommand (search by default): "
	fmt.Printf(searchMsg)

	// we'll read out std input via goroutine
	ch := make(chan string)
	go readInputChannel(ch)

	// read stdin and search for DB record
	patCopy, err := regexp.Compile(`cp [0-9]+`)
	if err != nil {
		log.Fatal(err)
	}
	patRemove, err := regexp.Compile(`rm [0-9]+`)
	if err != nil {
		log.Fatal(err)
	}
	patAdd, err := regexp.Compile(`add [login|card|note]`)
	if err != nil {
		log.Fatal(err)
	}
	patTimeout, err := regexp.Compile(`timeout [0-9]+`)
	if err != nil {
		log.Fatal(err)
	}
	dbRecords, err := readDB(db)
	if err != nil {
		log.Fatal(err)
	}

	// main loop
	for {
		select {
		case input := <-ch:
			input = strings.Replace(input, "\n", "", -1)
			if matched := patCopy.MatchString(input); matched {
				clipboardCopy(input, dbRecords)
			} else if matched := patRemove.MatchString(input); matched {
				removeRecord(input, dbRecords)
			} else if matched := patAdd.MatchString(input); matched {
				addRecord(input, dbRecords)
			} else if matched := patTimeout.MatchString(input); matched {
				vvv := strings.Trim(strings.Replace(input, "timeout ", "", -1), " ")
				if val, err := strconv.Atoi(vvv); err == nil {
					timeout = time.Duration(val)
					fmt.Printf("New DB timeout is set to %d seconds", timeout)
				}
				addRecord(input, dbRecords)
			} else {
				search(input, dbRecords)
			}
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

// helper function to copy to clipboard db record attribute
func clipboardCopy(input string, records *DBRecords) {
	// the input here is cp <ID> attribute
	arr := strings.Split(input, " ")
	if len(arr) < 2 {
		log.Printf("WARNING: unable to parse command '%s'", input)
		return
	}
	dbRecords := *records
	rid, err := strconv.Atoi(arr[1])
	if err != nil {
		log.Println("Unable to get record ID", err)
		return
	}
	attr := "password"
	if len(arr) == 3 {
		attr = strings.ToLower(arr[2])
	}
	var val string
	if entry, ok := dbRecords[rid]; ok {
		if attr == "password" {
			val = entry.GetPassword()
		} else if attr == "title" {
			val = getValue(entry, "Title")
		} else if attr == "username" {
			val = getValue(entry, "UserName")
		} else if attr == "url" {
			val = getValue(entry, "URL")
		} else if attr == "notes" {
			val = getValue(entry, "Notes")
		}
		if val != "" {
			if err := clipboard.WriteAll(string(val)); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s copied to clipboard\n", attr)
		}
	}
}
func removeRecord(input string, dbRecords *DBRecords) {
}
func addRecord(input string, dbRecords *DBRecords) {
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

// helper function to read db records
func readDB(db *gokeepasslib.Database) (*DBRecords, error) {
	records := make(DBRecords)

	for _, top := range db.Content.Root.Groups {
		if top.Name == "NewDatabase" {
			msg := "wrong password or empty database"
			return nil, errors.New(msg)
		}
		for _, groups := range top.Groups {
			rid := 0
			for _, entry := range groups.Entries {
				records[rid] = entry
				rid += 1
			}
		}
	}
	return &records, nil
}

// helper function to search for given input
func search(input string, records *DBRecords) {
	keys := []string{"UserName", "URL", "Notes"}
	pat := regexp.MustCompile(input)
	for rid, entry := range *records {
		if input == entry.GetTitle() || input == entry.Tags {
			printRecord(rid, entry)
		} else {
			for _, k := range keys {
				val := getValue(entry, k)
				if pat.MatchString(val) {
					printRecord(rid, entry)
				}
			}
		}
	}
}

// helper function to print record
func printRecord(pid int, entry gokeepasslib.Entry) {
	fmt.Printf("---\n")
	fmt.Printf("Record   %d\n", pid)
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
