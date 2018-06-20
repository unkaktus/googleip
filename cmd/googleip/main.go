package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/nogoegst/googleip"
)

func main() {
	var front = flag.String("f", "www.google.com", "front domain to use")
	flag.Parse()
	t := http.DefaultTransport
	ip, err := googleip.GetIP(t, *front)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", ip)
}
