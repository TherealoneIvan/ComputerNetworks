package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

// CommandType describes the different types of an Command.
type CommandType int

// The differents types of an Command
const (
	List CommandType = iota
	ChangeDir
	MakeDir
	Retr
	Stor
	DelFile
	Exit
	Undefined
)

// Command represents FTP command
type Command struct {
	Value string
	Path  string
	Arg   string
	Type  CommandType
}

// Execute : Выполнение команды
func (c Command) Execute(conn *ftp.ServerConn) {
	switch c.Type {
	case List:
		entries, err := conn.List(c.Arg)
		if err != nil {
			fmt.Println(err)
			return
		}
		for i, entry := range entries {
			fmt.Printf("%d: \"%s\" %d bytes [%s]\n", i, entry.Name, entry.Size, EntryTypeToString(entry))
		}
	case ChangeDir:
		err := conn.ChangeDir(c.Arg)
		if err != nil {
			fmt.Println(err)
		}
	case MakeDir:
		err := conn.MakeDir(c.Arg)
		if err != nil {
			fmt.Println(err)
		}
	case Retr:
		data, err := conn.Retr(c.Arg)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer data.Close()

		file, err := os.Create(fmt.Sprintf("%s/%d.zip", c.Path, time.Now().Unix()))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, data)
		if err != nil {
			fmt.Println(err)
		}
	case Stor:
		file, err := os.Open(c.Arg)
		if err != nil {
			fmt.Println(err)
			return
		}
		reader := bufio.NewReader(file)
		err = conn.Stor(file.Name(), reader)
		if err != nil {
			fmt.Println(err)
		}
	case DelFile:
		err := conn.Delete(c.Arg)
		if err != nil {
			fmt.Println(err)
		}
	case Exit:
		if err := conn.Quit(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("Завершение программы...")
	case Undefined:
		fmt.Println("Undefined command")
	}
}

// EntryTypeToString : Преобразует тип вхождения элемента ftp сервера к строке
func EntryTypeToString(e *ftp.Entry) string {
	switch e.Type {
	case ftp.EntryTypeFile:
		return "файл"
	case ftp.EntryTypeFolder:
		return "папка"
	case ftp.EntryTypeLink:
		return "ссылка"
	}
	return ""
}

func main() {
	var (
		host     string
		path     string
		port     int
		login    string
		password string
		command  Command
	)
	re := regexp.MustCompile(`\s+`)

	// установить указатели значений аргументов командной строки по ключам
	// на глобальные переменные аутентификации и инициализации сервера
	flag.StringVar(&host, "h", "185.20.227.83", "ftp server address")
	flag.StringVar(&path, "path", "/home/iu9_32_17/ftpclientroot", "path to files")
	flag.IntVar(&port, "port", 2017, "port of ftp server")

	// распарсить командную строку и инициализировать переменные значениями аргументов
	flag.Parse()

	// для отладки
	fmt.Println(host)
	fmt.Println(path)
	fmt.Println(port)

	c, err := ftp.Dial(fmt.Sprintf("%s:%d", host, port), ftp.DialWithTimeout(5*time.Second)) // соединение с ftp сервером
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Введите логин")
	fmt.Scanf("%s\n", &login)
	fmt.Println("Введите пароль")
	fmt.Scanf("%s\n", &password)

	err = c.Login(login, password) // авторизация на сервере
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Вы успешно соединились с FTP сервером!")

	scanner := bufio.NewScanner(os.Stdin)

	for command.Type != Exit {
		dir, err := c.CurrentDir()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Текущая директория: %s\n", dir)

		scanner.Scan()
		s := scanner.Text()

		command.Value = s
		command.Path = path

		re.ReplaceAllString(s, " ")
		data := strings.Split(s, " ")

		switch data[0] {
		case "ls":
			command.Arg = dir
			command.Type = List
		case "cd":
			if len(data) < 2 {
				fmt.Println(errors.New("err: missing Args"))
				return
			}
			command.Arg = data[1]
			command.Type = ChangeDir
		case "mkdir":
			if len(data) < 2 {
				fmt.Println(errors.New("err: missing Args"))
				return
			}
			command.Arg = data[1]
			command.Type = MakeDir
		case "retr":
			if len(data) < 2 {
				fmt.Println(errors.New("err: missing Args"))
				return
			}
			command.Arg = data[1]
			command.Type = Retr
		case "stor":
			if len(data) < 2 {
				fmt.Println(errors.New("err: missing Args"))
				return
			}
			command.Arg = data[1]
			command.Type = Stor
		case "rm":
			if len(data) < 2 {
				fmt.Println(errors.New("err: missing Args"))
				return
			}
			command.Arg = data[1]
			command.Type = DelFile
		case "exit":
			command.Type = Exit
		default:
			command.Type = Undefined
		}

		command.Execute(c)
	}
}
