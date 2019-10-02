package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
)

func main() {
	var host string
	flag.StringVar(&host, "host", "www.google.com", "host to trace")
	out, err := exec.Command("tracert", host).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
}
