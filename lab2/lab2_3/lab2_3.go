package main

import (
	"log"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		term := terminal.NewTerminal(s, "> ")
		for {
			line, err := term.ReadLine()
			if err != nil {
				break
			}
			response := line
			log.Println(line)
			if response != "" {
				term.Write(append([]byte(response), '\n'))
			}
		}
		log.Println("terminal closed")
	})

	log.Println("starting ssh server on port 2223...")
	log.Fatal(ssh.ListenAndServe(":2223", nil))
}
