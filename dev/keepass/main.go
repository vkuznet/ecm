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
	wrappers "github.com/tobischo/gokeepasslib/v3/wrappers"
	"golang.org/x/crypto/ssh/terminal"
)

// DBRecords defines map of DB records
type DBRecords map[int]gokeepasslib.Entry

// Record represent record map
type Record map[string]string

// keep track if user requested password field
var inputPwd bool

// helper function to print commands usage
func cmdUsage(dbPath string) {
	if dbPath != "" {
		info, err := os.Stat(dbPath)
		if err == nil {
			fmt.Println()
			fmt.Println("Database            : ", dbPath)
			fmt.Println("Size                : ", sizeFormat(info.Size()))
			fmt.Println("Modification time   : ", info.ModTime())
			fmt.Println()
		}
	}
	fmt.Println("Commands within DB")
	fmt.Println("cp <ID> <attribute> # to copy record ID attribute to cpilboard")
	fmt.Println("rm <ID>             # to remove record ID from database")
	fmt.Println("add <key>           # to add specific record key")
	fmt.Println("timeout <int>       # set timeout interval in seconds")
}

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
		cmdUsage("")
	}
	flag.Parse()

	file, err := os.Open(kpath)
	if err != nil {
		log.Fatal(err)
	}

	db := gokeepasslib.NewDatabase()
	pwd := readPassword("db password: ")
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
	dbRecords, err := readDB(db)
	if err != nil {
		log.Fatal(err)
	}

	// proceed with db records
	cmdUsage(kpath)
	var names []string
	for _, g := range db.Content.Root.Groups {
		names = append(names, g.Name)
	}
	fmt.Printf("Welcome to %s", strings.Join(names, ","))

	inputMsg := "\ndb # "
	inputMsgOrig := inputMsg
	inputPwd = false
	fmt.Printf(inputMsg)

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
	patAdd, err := regexp.Compile(`add [a-zA-Z]+`)
	if err != nil {
		log.Fatal(err)
	}
	//     patSave, err := regexp.Compile(`save record`)
	//     if err != nil {
	//         log.Fatal(err)
	//     }
	patTimeout, err := regexp.Compile(`timeout [0-9]+`)
	if err != nil {
		log.Fatal(err)
	}

	// main loop
	var rec Record
	rec = nil
	collectKey := ""
	for {
		select {
		case input := <-ch:
			input = strings.Replace(input, "\n", "", -1)
			if input == "save record" {
				inputPwd = false
				kind := "Record"
				if v, ok := rec["Login"]; ok {
					rec["UserName"] = v
					kind = "Login"
				} else if _, ok := rec["Card"]; ok {
					kind = "Card"
				} else if _, ok := rec["Note"]; ok {
					kind = "Note"
				}
				saveRecord(kind, rec, db)
				rec = nil
				collectKey = ""
				inputMsg = inputMsgOrig
			} else if input == "exit" || input == "quit" {
				os.Exit(0)
			} else if strings.HasPrefix(input, "WARNING") {
				inputPwd = false
				collectKey = ""
				fmt.Println(input)
				inputMsg = inputMsgOrig
			} else if collectKey != "" {
				inputPwd = false
				collectKey = ""
				rec[collectKey] = input
				inputMsg = inputMsgOrig
			} else if matched := patCopy.MatchString(input); matched {
				inputPwd = false
				clipboardCopy(input, dbRecords)
				inputMsg = inputMsgOrig
			} else if matched := patRemove.MatchString(input); matched {
				inputPwd = false
				removeRecord(input, dbRecords)
				inputMsg = inputMsgOrig
			} else if matched := patAdd.MatchString(input); matched {
				if rec == nil {
					rec = make(map[string]string)
				}
				collectKey = strings.Replace(input, "add ", "", -1)
				if strings.ToLower(collectKey) == "password" {
					inputPwd = true
				} else {
					inputPwd = false
				}
				inputMsg = fmt.Sprintf("%s value: ", collectKey)
			} else if matched := patTimeout.MatchString(input); matched {
				inputPwd = false
				vvv := strings.Trim(strings.Replace(input, "timeout ", "", -1), " ")
				if val, err := strconv.Atoi(vvv); err == nil {
					timeout = time.Duration(val)
					fmt.Printf("New DB timeout is set to %d seconds", timeout)
				}
				inputMsg = inputMsgOrig
			} else {
				inputPwd = false
				search(input, dbRecords)
				inputMsg = inputMsgOrig
			}
			time0 = time.Now()
			fmt.Printf(inputMsg)
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

// helper function to remove record from the database
func removeRecord(input string, dbRecords *DBRecords) {
}

// helper function to make entry db value
func mkValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{Key: key, Value: gokeepasslib.V{Content: value}}
}

// helper function to make protected entry db value
func mkProtectedValue(key string, value string) gokeepasslib.ValueData {
	return gokeepasslib.ValueData{
		Key:   key,
		Value: gokeepasslib.V{Content: value, Protected: wrappers.NewBoolWrapper(true)},
	}
}

// helper function to save record to the database
func saveRecord(kind string, rec Record, db *gokeepasslib.Database) {
	group := gokeepasslib.NewGroup()
	group.Name = kind
	entry := gokeepasslib.NewEntry()
	for key, val := range rec {
		attr := strings.ToLower(key)
		if attr == "password" {
			entry.Values = append(entry.Values, mkProtectedValue("Password", val))
		} else {
			entry.Values = append(entry.Values, mkValue(key, val))
		}
	}
	group.Entries = append(group.Entries, entry)
	log.Printf("added %s db entry: %+v", kind, entry)
	// write group entries to DB
	// https://github.com/tobischo/gokeepasslib/blob/master/examples/writing/example-writing.go

	/*

		// now create the database containing the root group
		db := &gokeepasslib.Database{
			Header:      gokeepasslib.NewHeader(),
			Credentials: gokeepasslib.NewPasswordCredentials(masterPassword),
			Content: &gokeepasslib.DBContent{
				Meta: gokeepasslib.NewMetaData(),
				Root: &gokeepasslib.RootData{
					Groups: []gokeepasslib.Group{rootGroup},
				},
			},
		}

		// Lock entries using stream cipher
		db.LockProtectedEntries()

		// and encode it into the file
		keepassEncoder := gokeepasslib.NewEncoder(file)
		if err := keepassEncoder.Encode(db); err != nil {
			panic(err)
		}

		log.Printf("Wrote kdbx file: %s", filename)
	*/
	log.Println("Updated kdbx file")
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
	var val string
	var err error
	for {
		//         fmt.Print(*msg)
		if inputPwd == true {
			val = readPassword("")
			// read password again to match it
			if v := readPassword("repeat password: "); v != val {
				ch <- fmt.Sprintf("WARNING: password match failed, will discard it ...")
				continue
			}
		} else {
			val, err = readInput()
			if err != nil {
				ch <- fmt.Sprintf("WARNING: wrong input %v", err)
				continue
			}
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
func readPassword(msg string) string {
	if msg != "" {
		fmt.Print(msg)
	}
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

// helper function to convert size into human readable form
func sizeFormat(val int64) string {
	size := float64(val)
	base := 1000. // CMS convert is to use power of 10
	xlist := []string{"", "KB", "MB", "GB", "TB", "PB"}
	for _, vvv := range xlist {
		if size < base {
			return fmt.Sprintf("%v (%3.1f%s)", val, size, vvv)
		}
		size = size / base
	}
	return fmt.Sprintf("%v (%3.1f%s)", val, size, xlist[len(xlist)])
}
