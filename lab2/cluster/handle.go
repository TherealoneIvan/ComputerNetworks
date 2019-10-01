package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Server gives info about server to connect
type Server struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func main() {
	var path, cmd string
	flag.StringVar(&path, "path", "cluster.json", "path to cluster info file")
	flag.StringVar(&cmd, "cmd", "pwd", "command to run")
	flag.Parse()

	var servers []Server

	// открытие JSON файла со списком серверов
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// считывание информации из файла в массив байтов
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// десериализация JSON строки в экземпляр структуры
	json.Unmarshal(byteValue, &servers)

	resChannel := make(chan string, len(servers))
	timeout := time.After(5 * time.Second)

	for i := 0; i < len(servers); i++ {
		go func(j int) {
			config := &ssh.ClientConfig{
				User: servers[j].User,
				Auth: []ssh.AuthMethod{
					ssh.Password(servers[j].Password)},
				HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
					return nil
				}}

			conn, _ := ssh.Dial("tcp", fmt.Sprintf("%s:%d", servers[j].Host, servers[j].Port), config)
			session, _ := conn.NewSession()
			defer session.Close()

			var stdoutBuf bytes.Buffer
			session.Stdout = &stdoutBuf
			session.Run(cmd)

			resChannel <- fmt.Sprintf("%s\n\n%s\n", servers[j].Host, stdoutBuf.String())
		}(i)
	}

	for i := 0; i < len(servers); i++ {
		select {
		case res := <-resChannel:
			fmt.Print(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
}
