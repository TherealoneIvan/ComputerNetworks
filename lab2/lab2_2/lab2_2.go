package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/gliderlabs/ssh"
)

func main() {
	var port int

	flag.IntVar(&port, "port", 2217, "port of ssh server")

	flag.Parse()

	server := &ssh.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: func(s ssh.Session) {
			io.WriteString(s, fmt.Sprintf("You've been connected to %s\n", s.LocalAddr().String()))
			for {
				text, err := bufio.NewReader(s).ReadString('\n')
				if err != nil {
					fmt.Println("GetLines: " + err.Error())
					err = s.Exit(-1)
					if err != nil {
						log.Fatal(err)
					}
				}
				fmt.Println(text)
				if text == "exit" {
					err = s.Exit(0)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		},
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			return ctx.User() == "iu9_student" && password == "BMSTU_the_best"
		},
	}

	log.Println(fmt.Sprintf("starting ssh server on port %d...", port))
	log.Fatal(server.ListenAndServe())
}
