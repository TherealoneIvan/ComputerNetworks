package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"           // пакет для форматированного ввода вывода
	"html/template" // пакет для шаблонизации HTML документов
	"io"            // пакет для работы с вводом/выводом
	"log"           // пакет для логирования
	"net/http"      // пакет для поддержки HTTP протокола
	"sort"          // пакет для сортировки
	"strconv"       // пакет для конвертации строк  в другой тип
	"strings"       // пакет для работы с  UTF-8 строками
	"time"          // пакет для работы со временем

	"github.com/jlaffaye/ftp" //пакет для написания FTP клиента
)

// FTPIdentity represents data for auth
type FTPIdentity struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

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
func EntriesByType(path string, conn *ftp.ServerConn, etype ftp.EntryType) ([]EntryExtended, error) {
	entries, err := conn.List(path)
	if err != nil {
		return nil, err
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
	return filteredEntries, nil
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
func FTPAuth(path string, w http.ResponseWriter, r *http.Request) *ftp.ServerConn {
	var identity FTPIdentity

	cookie, err := r.Cookie("identity")
	if err != nil {
		ErrorHandling(err, path, w)
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		ErrorHandling(err, path, w)
		return nil
	}

	err = json.Unmarshal(data, &identity)
	if err != nil {
		ErrorHandling(err, path, w)
		return nil
	}

	c, err := ftp.Dial(fmt.Sprintf("%s:%d", identity.Host, identity.Port), ftp.DialWithTimeout(5*time.Second)) // соединение с ftp сервером
	if err != nil {
		ErrorHandling(err, path, w)
		return nil
	}

	err = c.Login(identity.Login, identity.Password) // авторизация на сервере
	if err != nil {
		ErrorHandling(err, path, w)
		return nil
	}

	return c
}

// AuthRouterHandler : Запись данных для авторизации на FTP сервере
func AuthRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	i, err := strconv.Atoi(r.Form["port"][0])
	if err != nil {
		ErrorHandling(err, "/", w)
		return
	}

	identity := FTPIdentity{
		Host:     r.Form["host"][0],
		Port:     i,
		Login:    r.Form["login"][0],
		Password: r.Form["password"][0],
	}

	jsonBytes, _ := json.Marshal(identity)

	http.SetCookie(w, &http.Cookie{
		Name:    "identity",
		Value:   base64.StdEncoding.EncodeToString(jsonBytes),
		Expires: time.Now().Add(30 * time.Minute),
	})

	http.Redirect(w, r, "/client/", 301)
}

// ClientRouterHandler : Обработчик запросов ко клиенту
func ClientRouterHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/client", "", 1)
	n := len(path)
	if n > 1 && path[n-1] == '/' {
		path = path[:n-1]
	}

	c := FTPAuth(path, w, r)

	var entries []EntryExtended
	files, err := EntriesByType(path, c, ftp.EntryTypeFolder) // папки
	if err != nil {
		ErrorHandling(err, path, w)
		return
	}
	folders, err := EntriesByType(path, c, ftp.EntryTypeFile) // файлы
	links, err := EntriesByType(path, c, ftp.EntryTypeLink)   // ссылки

	entries = append(entries, files...)
	entries = append(entries, folders...)
	entries = append(entries, links...)

	if err := c.Quit(); err != nil {
		ErrorHandling(err, path, w)
		return
	}

	tmplt := template.Must(template.ParseFiles("static/index.html"))
	data := EntriesPageData{
		EntriesList: entries,
		CurrentDir:  path}
	tmplt.Execute(w, data) // отправляем данные на клиентскую сторону
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

	c := FTPAuth(current, w, r)

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

	c := FTPAuth(current, w, r)

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

	c := FTPAuth(current, w, r)

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

	c := FTPAuth(current, w, r)

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
