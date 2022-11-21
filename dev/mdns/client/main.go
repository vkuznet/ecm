package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/hashicorp/mdns"
)

func main() {

	serviceTag := "_mdns._tcp"
	if len(os.Args) > 1 {
		serviceTag = os.Args[1]
	}

	// make verbose logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Make a channel for results and start listening
	ch := make(chan *mdns.ServiceEntry, 8)
	defer close(ch)

	go func() {
		for entry := range ch {
			//             fmt.Printf("Got new entry: %v\n", entry)
			fmt.Printf("Discovered name=%s at http://%s:%d\n", entry.Name, entry.AddrV4, entry.Port)
		}
	}()

	params := mdns.DefaultParams(serviceTag)
	params.DisableIPv6 = true
	params.Entries = ch
	err := mdns.Query(params)
	if err != nil {
		log.Fatal(err)
	}

	// Start the lookups
	/*
		err := mdns.Lookup(serviceTag, ch)
		if err != nil {
			fmt.Println(err)
		}
	*/

	wait()
}

func wait() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}
