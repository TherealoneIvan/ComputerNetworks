package main

import (
	"fmt"           // пакет для форматированного ввода вывода
	"html/template" // пакет для шаблонизации HTML документов
	"log"           // пакет для логирования
	"net/http"      // пакет для поддержки HTTP протокола
	"sort"          // пакет для сортировки
	"strings"       // пакет для работы с  UTF-8 строками
	"time"          // пакет для работы со временем

	"github.com/jlaffaye/ftp" //пакет для написания FTP клиента
)

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

// HomeRouterHandler : router
func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	c, err := ftp.Dial("students.yss.su:21", ftp.DialWithTimeout(5*time.Second)) // соединение с ftp сервером
	if err != nil {
		fmt.Println(err)
	}

	err = c.Login("ftpiu8", "3Ru7yOTA") // авторизация на сервере
	if err != nil {
		fmt.Println(err)
	}

	r.ParseForm() // анализ аргументов

	path := strings.Replace(r.URL.Path, "/client", "", 1)
	n := len(path)
	if n > 1 && path[n-1] == '/' {
		path = path[:n-1]
	}

	var entries []EntryExtended
	entries = append(entries, EntriesByType(path, c, ftp.EntryTypeFolder)...) // папки
	entries = append(entries, EntriesByType(path, c, ftp.EntryTypeFile)...)   // файлы
	entries = append(entries, EntriesByType(path, c, ftp.EntryTypeLink)...)   // ссылки

	if err := c.Quit(); err != nil {
		fmt.Println(err)
	}

	if len(entries) == 0 {
		fmt.Fprintf(w, "No such file or directory")
	} else {
		tmplt := template.Must(template.ParseFiles("static/index.html"))
		data := EntriesPageData{
			EntriesList: entries,
			CurrentDir:  path}
		tmplt.Execute(w, data) // отправляем данные на клиентскую сторону
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/client", 301)
	})
	http.HandleFunc("/client/", HomeRouterHandler) // установим роутер
	err := http.ListenAndServe(":8017", nil)       // задаём слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
