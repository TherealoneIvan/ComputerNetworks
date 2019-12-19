package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func factorial(n int64) (result int64) {
	if n > 0 {
		result = n * factorial(n-1)
		return result
	}
	return 1
}

// требуется только ниже для обработки примера

func main() {
	var port int
	flag.IntVar(&port, "port", 8002, "port of tcp server")

	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))

	// Открываем порт
	conn, _ := ln.Accept()

	// Запускаем цикл
	for {
		// Будем прослушивать все сообщения разделенные \n
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// Распечатываем полученое сообщение
		fmt.Print("Message Received:", string(message))
		i64, err := strconv.ParseInt(strings.TrimRight(message, "\r\n"), 10, 64)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
		}
		// Отправить новую строку обратно клиенту
		conn.Write([]byte(strconv.FormatInt(factorial(i64), 10) + "\n"))
	}
}
