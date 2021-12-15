package main

import (
	"fmt"
	"syscall/js"
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
		input := args[0].String()
		// perform decoding of the input
		output := fmt.Sprintf("GPM: %s", input)
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
