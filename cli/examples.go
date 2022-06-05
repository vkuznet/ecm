package main

import "fmt"

// helper function to print examples
func ecmExamples() {
	fmt.Println("ECM examples:")
	fmt.Println("# get list of records from the vault")
	fmt.Println("./ecm")
	fmt.Println("# or explicitly specify vault location")
	fmt.Println("./ecm -vault /path/to/vault")
	fmt.Println("")
	fmt.Println("# get vault info")
	fmt.Println("./ecm -info")
	fmt.Println("")
	fmt.Println("# get info about single vault record (and its password will be copied to clipboard)")
	fmt.Println("./ecm -rid cc1ee1e4-183c-423f-9ce1-62f26287441b")
	fmt.Println("")
	fmt.Println("# edit given vault record")
	fmt.Println("./ecm -rid cc1ee1e4-183c-423f-9ce1-62f26287441b")
	fmt.Println("")
	fmt.Println("# decrypt file")
	fmt.Println("./ecm -decrypt ~/.ecm/Primary/cc1ee1e4-183c-423f-9ce1-62f26287441b")
	fmt.Println("")
	fmt.Println("# encrypt file, new recrod will be created in a vault")
	fmt.Println("./ecm -decrypt some-file")
	fmt.Println("")
	fmt.Println("# sync vault from some local fila system location")
	fmt.Println("./ecm -sync=file:///tmp/ecm")
	fmt.Println("# or sync from dropbox")
	fmt.Println("# storage should be configure via rclone")
	fmt.Println("# configuration file usually located at $HOME/.config/rclone/rclone.conf")
	fmt.Println("./ecm -sync=dropbox:ECM")
}
