package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
				case "echo":
					if len(command) < 2 {
						io.WriteString(s, fmt.Sprintf("%s: missing arg\n", command[0]))
					} else {
						io.WriteString(s, fmt.Sprintf("%s\n", command[1]))
					}
				case "ls":
					files, err := ioutil.ReadDir("./")
					if err != nil {
						fmt.Println(err)
						io.WriteString(s, err.Error())
					}

					for _, f := range files {
						io.WriteString(s, fmt.Sprintf("%s\n", f.Name()))
					}
				case "mkdir":
					if len(command) < 2 {
						io.WriteString(s, fmt.Sprintf("%s: missing arg\n", command[0]))
					} else {
						err := os.Mkdir(command[1], os.ModeDir)
						if err != nil {
							fmt.Println(err)
							io.WriteString(s, err.Error())
						}
					}
				case "rmdir":
					if len(command) < 2 {
						io.WriteString(s, fmt.Sprintf("%s: missing arg\n", command[0]))
					} else {
						err := os.Remove(command[1])
						if err != nil {
							fmt.Println(err)
							io.WriteString(s, err.Error())
						}
					}
				case "exit":
					break loop
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
