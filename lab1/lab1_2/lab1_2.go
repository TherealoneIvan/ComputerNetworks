package main

import (
	"flag"
	"fmt"

	filedriver "gitea.com/goftp/file-driver"
	"goftp.io/server"
)

func main() {
	var (
		host     string
		path     string
		port     int
		login    string
		password string
	)

	// установить указатели значений аргументов командной строки по ключам
	// на глобальные переменные аутентификации и инициализации сервера
	flag.StringVar(&host, "h", "185.20.227.83", "ftp server address")
	flag.StringVar(&path, "path", "/home/iu9_32_17/ftproot", "ftp server root folder")
	flag.IntVar(&port, "port", 2117, "port of ftp server")
	flag.StringVar(&login, "l", "login", "login for ftp auth")
	flag.StringVar(&password, "p", "p4$$w0rd", "password for ftp auth")

	// распарсить командную строку и инициализировать переменные значениями аргументов
	flag.Parse()

	// для отладки
	fmt.Println(host)
	fmt.Println(path)
	fmt.Println(port)
	fmt.Println(login)
	fmt.Println(password)

	factory := &filedriver.FileDriverFactory{
		RootPath: path,
		Perm:     server.NewSimplePerm("root", "root")}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     port,
		Hostname: host,
		Auth:     &server.SimpleAuth{Name: login, Password: password}}

	server := server.NewServer(opts)
	server.ListenAndServe()
}
