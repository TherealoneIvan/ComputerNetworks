package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gliderlabs/ssh"
)

func parseCommand(s string) []string {
	s = strings.TrimRight(s, "\n")
	re := regexp.MustCompile(`\s+`)
	re.ReplaceAllString(s, " ")
	data := strings.Split(s, " ")
	return data
}

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
					break
				}

				fmt.Println(text)
				command := parseCommand(text)

				out, err := exec.Command(command[0], command[1:]...).Output()
				if err != nil {
					log.Fatal(err)
				}
				io.WriteString(s, string(out))
			}
			err := s.Exit(0)
			if err != nil {
				fmt.Println(err)
			}
		},
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			return ctx.User() == "iu9_student" && password == "BMSTU_the_best"
		},
	}

	log.Println(fmt.Sprintf("starting ssh server on port %d...", port))
	log.Fatal(server.ListenAndServe())
}
