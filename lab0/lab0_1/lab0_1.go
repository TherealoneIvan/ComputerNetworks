package main

import (
	"fmt"      // пакет для форматированного ввода вывода
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола
	"strings"  // пакет для работы с  UTF-8 строками
)

// HomeRouterHandler : router
func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html") // возвращать HTML документ
	r.ParseForm()                               // анализ аргументов
	var response strings.Builder
	response.WriteString("<ul>")
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, " "))
		response.WriteString(fmt.Sprintf("<li>%s : %s</li>", k, strings.Join(v, " ")))
	}
	response.WriteString("</ul>")
	if r.URL.Path != "/info" {
		response.WriteString("<iframe src=\"http://localhost:8017/info?lang=go&version=1.12.7&os=windows\">Link</iframe>")
	}
	fmt.Fprintf(w, response.String()) // отправляем данные на клиентскую сторону
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)  // установим роутер
	err := http.ListenAndServe(":8017", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
