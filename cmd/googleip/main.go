package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/nogoegst/frontier"
	"github.com/nogoegst/googleip"
)

func main() {
	var front = flag.String("f", "www.google.com", "front domain to use")
	var addr = flag.String("a", "", "address to use")
	flag.Parse()
	t := http.DefaultTransport
	ip, err := googleip.GetIP(frontier.New(t, *front, *addr))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", ip)
}
