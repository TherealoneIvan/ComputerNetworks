package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8002, "port of tcp server")

	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	for {
		// Чтение входных данных от stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("number to send: ")
		text, _ := reader.ReadString('\n')
		if strings.TrimRight(text, "\r\n") == "exit" {
			conn.Close()
			os.Exit(0)
		}
		// Отправляем в socket
		fmt.Fprintf(conn, text+"\n")
		// Прослушиваем ответ
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
