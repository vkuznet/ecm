package main

import (
	"net/http"
)

func main() {
	panic(http.ListenAndServe(":9090", http.FileServer(http.Dir("."))))
}
