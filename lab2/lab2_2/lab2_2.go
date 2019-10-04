package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
		loop:
			for {
				text, err := bufio.NewReader(s).ReadString('\n')
				if err != nil {
					fmt.Println("GetLines: " + err.Error())
					break
				}
				fmt.Println(text)

				command := parseCommand(text)

				switch command[0] {
				case "exit":
					break loop
				case "cd":
					if len(command) < 2 {
						home := os.Getenv("HOME")
						os.Chdir(home)
					} else {
						err := os.Chdir(command[1])
						if err != nil {
							io.WriteString(s, err.Error())
						}
					}
				default:
					out, err := exec.Command(command[0], command[1:]...).Output()
					if err != nil {
						io.WriteString(s, err.Error())
					}
					io.WriteString(s, fmt.Sprintf("%s\n", string(out)))
				}
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
