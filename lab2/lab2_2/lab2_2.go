package main

import (
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
			io.WriteString(s, fmt.Sprintf("You've been connected to %s:%d\n", s.LocalAddr().String(), port))
		},
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			return ctx.User() == "iu9_student" && password == "BMSTU_the_best"
		},
	}

	log.Println(fmt.Sprintf("starting ssh server on port %d...", port))
	log.Fatal(server.ListenAndServe())
}
