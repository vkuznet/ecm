package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"syscall/js"

	"github.com/vkuznet/gpm/crypt"
)

// main function sets JS "gpmDecode" function to call "decodeWrapper" Go counterpart
func main() {
	js.Global().Set("gpmDecode", decodeWrapper())
	<-make(chan bool)
}

// decodeWraper function performs business logic, i.e.
// it decodes given input obtained from JS upstream code
func decodeWrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return wrap("", "Not enough arguments")
		}
		// here input should be vault file name, e.g.
		// "/Users/vk/.gpm/Primary/acb8a9f7-6140-42d2-bb32-f730f7ab572f.aes"
		input := args[0].String()
		// perform decoding input file
		passphrase := "test"
		cipher := "aes"
		arr := strings.Split(input, ".")
		if len(arr) > 0 {
			cipher = arr[1]
		}
		data, err := ioutil.ReadFile(input)
		if err != nil {
			wrap("", err.Error())
		}
		rec, err := crypt.Decrypt(data, passphrase, cipher)
		if err != nil {
			wrap("", err.Error())
		}
		output := fmt.Sprintf("GPM: %s, data=%s", input, string(rec))
		return wrap(output, "")
	})
}

// wrap function provide output to JS upstream code in form of the dict
func wrap(record string, err string) map[string]interface{} {
	return map[string]interface{}{
		"error":  err,
		"record": record,
	}
}
