package main

import (
	"errors" // пакет для работы с ошибками
	"strconv"

	// пакет для работы с аргументами командной строки
	"fmt"           // пакет для форматированного ввода вывода
	"html/template" // пакет для шаблонизации HTML документов
	"io"            // пакет для работы с вводом/выводом
	"log"           // пакет для логирования
	"net/http"      // пакет для поддержки HTTP протокола
	"sort"          // пакет для сортировки
	"strings"       // пакет для работы с  UTF-8 строками
	"time"          // пакет для работы со временем

	"github.com/jlaffaye/ftp" //пакет для написания FTP клиента
)

// данные авторизации
var host string // адрес ftp сервера
var port int    // порт ftp сервера
var login string
var password string

// EntryExtended extension of Entry type
type EntryExtended struct {
	Content      *ftp.Entry
	AbsolutePath string
}

// EntriesPageData data for html templating
type EntriesPageData struct {
	EntriesList []EntryExtended
	CurrentDir  string
}

// ErrorPageData data for error page
type ErrorPageData struct {
	Message     string
	RedirectURL string
}

// EntryTypeToString : Преобразует тип вхождения элемента ftp сервера к строке
func (ee EntryExtended) EntryTypeToString() string {
	switch ee.Content.Type {
	case ftp.EntryTypeFile:
		return "файл"
	case ftp.EntryTypeFolder:
		return "папка"
	case ftp.EntryTypeLink:
		return "ссылка"
	}
	return ""
}

// EntriesByType : Возвращает список вхождений по заданному типу вхождения, отсортированный по алфавиту
func EntriesByType(path string, conn *ftp.ServerConn, etype ftp.EntryType) []EntryExtended {
	entries, err := conn.List(path)
	if err != nil {
		fmt.Println(err)
	}
	var filteredEntries []EntryExtended
	for _, entry := range entries {
		if entry.Type == etype {
			ee := EntryExtended{
				Content:      entry,
				AbsolutePath: strings.Join([]string{path, entry.Name}, "/")}
			filteredEntries = append(filteredEntries, ee)
		}
	}
	sort.Slice(filteredEntries, func(i, j int) bool {
		return strings.Compare(filteredEntries[i].Content.Name, filteredEntries[j].Content.Name) == -1
	})
	return filteredEntries
}

// ErrorHandling : Возврат клиенту информации о произошедшей ошибке
func ErrorHandling(e error, path string, w http.ResponseWriter) {
	tmplt := template.Must(template.ParseFiles("static/error.html"))
	data := ErrorPageData{
		Message:     e.Error(),
		RedirectURL: path}
	tmplt.Execute(w, data) // отправляем данные на клиентскую сторону
}

// FTPAuth : Авторизация на FTP сервере
func FTPAuth(path string, w http.ResponseWriter) *ftp.ServerConn {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", host, port), ftp.DialWithTimeout(5*time.Second)) // соединение с ftp сервером
	if err != nil {
		ErrorHandling(err, path, w)
	}

	err = c.Login(login, password) // авторизация на сервере
	if err != nil {
		ErrorHandling(err, path, w)
	}

	return c
}

// AuthRouterHandler : Запись данных для авторизации на FTP сервере
func AuthRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	host = r.Form["host"][0]
	i, err := strconv.Atoi(r.Form["port"][0])
	if err != nil {
		ErrorHandling(err, "/", w)
		return
	}
	port = i
	login = r.Form["login"][0]
	password = r.Form["password"][0]

	http.Redirect(w, r, "/client/", 301)
}

// ClientRouterHandler : Обработчик запросов ко клиенту
func ClientRouterHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/client", "", 1)
	n := len(path)
	if n > 1 && path[n-1] == '/' {
		path = path[:n-1]
	}

	c := FTPAuth(path, w)

	var entries []EntryExtended
	entries = append(entries, EntriesByType(path, c, ftp.EntryTypeFolder)...) // папки
	entries = append(entries, EntriesByType(path, c, ftp.EntryTypeFile)...)   // файлы
	entries = append(entries, EntriesByType(path, c, ftp.EntryTypeLink)...)   // ссылки

	if err := c.Quit(); err != nil {
		ErrorHandling(err, path, w)
		return
	}

	if len(entries) == 0 {
		ErrorHandling(errors.New("No such file or directory"), path, w)
	} else {
		tmplt := template.Must(template.ParseFiles("static/index.html"))
		data := EntriesPageData{
			EntriesList: entries,
			CurrentDir:  path}
		tmplt.Execute(w, data) // отправляем данные на клиентскую сторону
	}
}

// CreateRouterHandler : Обработчик запроса на создание директории
func CreateRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // анализ аргументов

	current := "" // текущая директория
	dir := ""     // имя создаваемой директории

	for k, v := range r.Form {
		if k == "curdir" {
			current = v[0]
		} else if k == "dirname" {
			dir = v[0]
		}
	}

	c := FTPAuth(current, w)

	err := c.ChangeDir(current)
	if err != nil {
		ErrorHandling(err, current, w)
		return
	}

	err = c.MakeDir(dir)
	if err != nil {
		ErrorHandling(err, current, w)
		return
	}

	if err := c.Quit(); err != nil {
		ErrorHandling(err, current, w)
		return
	}

	http.Redirect(w, r, strings.Join([]string{"client", current}, "/"), 301)
}

// UploadRouterHandler : Обработчик запроса на загрузку файла
func UploadRouterHandler(w http.ResponseWriter, r *http.Request) {
	var maxMB int64 = 10 // максимальный размер загружаемого файла в мегабайтах

	r.ParseMultipartForm(maxMB << 20) // анализ аргументов

	current := r.Form["curdir"][0] // директория загрузки
	file, handler, err := r.FormFile("uploadedfile")
	if err != nil {
		ErrorHandling(err, current, w)
		return
	}
	defer file.Close()

	c := FTPAuth(current, w)

	err = c.ChangeDir(current)
	if err != nil {
		ErrorHandling(err, current, w)
		return
	}

	err = c.Stor(handler.Filename, file)
	if err != nil {
		panic(err)
	}

	if err := c.Quit(); err != nil {
		ErrorHandling(err, current, w)
		return
	}

	http.Redirect(w, r, strings.Join([]string{"client", current}, "/"), 301)
}

// DownloadRouterHandler : Обработчик запроса на скачивание файла
func DownloadRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // анализ аргументов

	current := "" // текущая директория
	file := ""    // имя скачиваемого файла

	for k, v := range r.Form {
		if k == "curdir" {
			current = v[0]
		} else if k == "filename" {
			file = v[0]
		}
	}

	c := FTPAuth(current, w)

	err := c.ChangeDir(current)
	if err != nil {
		ErrorHandling(err, current, w)
		return
	}

	reader, err := c.Retr(file) // считывающий поток, содержащий байты файла
	if err != nil {
		ErrorHandling(err, current, w)
		return
	}

	// установка необходимых заголовков
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file))

	if err := c.Quit(); err != nil {
		ErrorHandling(err, current, w)
		return
	}

	io.Copy(w, reader) // отдать клиенту читающий поток без загрузки файла в память
}

// DeleteFileRouterHandler : Обработчик запроса на удаление файла
func DeleteFileRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // анализ аргументов

	current := "" // текущая директория
	file := ""    // имя удаляемого файла

	for k, v := range r.Form {
		if k == "curdir" {
			current = v[0]
		} else if k == "filename" {
			file = v[0]
		}
	}

	c := FTPAuth(current, w)

	err := c.ChangeDir(current)
	if err != nil {
		ErrorHandling(err, current, w)
		return
	}

	err = c.Delete(file)
	if err != nil {
		ErrorHandling(err, current, w)
		return
	}

	if err := c.Quit(); err != nil {
		ErrorHandling(err, current, w)
		return
	}

	http.Redirect(w, r, strings.Join([]string{"client", current}, "/"), 301)
}

func main() {
	// установка дефолтных значений глобальных переменных аутентификации
	host = "students.yss.su"
	port = 21
	login = "ftpiu8"
	password = "3Ru7yOTA"

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css")))) // доступ к стилям

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/auth.html")
	}) // возврат страница авторизации по запросу домашней страницы

	// установка обработчиков запросов
	http.HandleFunc("/auth", AuthRouterHandler)
	http.HandleFunc("/client/", ClientRouterHandler)
	http.HandleFunc("/create", CreateRouterHandler)
	http.HandleFunc("/upload", UploadRouterHandler)
	http.HandleFunc("/download", DownloadRouterHandler)
	http.HandleFunc("/deletefile", DeleteFileRouterHandler)

	err := http.ListenAndServe(":8017", nil) // задаём слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
