package main

import (
	"fmt"      // пакет для форматированного ввода вывода
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола
	"strings"  // пакет для работы с  UTF-8 строками

	"github.com/RealJK/rss-parser-go" //пакет для парсинга RSS каналов
)

// NewsRouterHandler : обработчик страницы выдачи новостей
func NewsRouterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") // возвращать HTML документ
	r.ParseForm()                                              // анализ аргументов

	var response strings.Builder // конструктор HTML-строки ответа

	response.WriteString("<div>")
	response.WriteString("<a href=\"/\">К списку каналов</a>")

	url := r.Form["url"][0] //ссылка на RSS канал
	rssObject, rssErr := rss.ParseRSS(url)
	if rssErr != nil {

		response.WriteString(fmt.Sprintf("<h1>%s</h1>", rssObject.Channel.Title))
		response.WriteString(fmt.Sprintf("<p>%s</p><br>", rssObject.Channel.Description))

		response.WriteString(fmt.Sprintf("<p>Количество новостей: %d</p><br>", len(rssObject.Channel.Items)))

		response.WriteString("<ul>")
		for v := range rssObject.Channel.Items {
			response.WriteString("<li><div>")
			item := rssObject.Channel.Items[v]

			response.WriteString(fmt.Sprintf("<h2>%s</h2>", item.Title))
			response.WriteString("<div>")
			response.WriteString(fmt.Sprintf("<img src=\"%s\" />", item.Enclosure.Url))
			if strings.Contains(item.Description, ">") {
				response.WriteString(fmt.Sprintf("<p>%s</p>", strings.Split(item.Description, ">")[1]))
			} else {
				response.WriteString(fmt.Sprintf("<p>%s</p>", item.Description))
			}
			response.WriteString(fmt.Sprintf("</div><form action=\"%s\" method=\"GET\">", item.Link))
			response.WriteString(fmt.Sprintf("<input type=\"submit\" value=\"К странице новости\"></form></div></li>"))
		}
		response.WriteString("</ul>")
	}

	response.WriteString("</div>")

	fmt.Fprintf(w, response.String()) // отправляем данные на клиентскую сторону
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static"))) // установим роутер home
	http.HandleFunc("/news", NewsRouterHandler)           // установим роутер news
	httpErr := http.ListenAndServe(":8017", nil)          // задаем слушать порт
	if httpErr != nil {
		log.Fatal("ListenAndServe: ", httpErr)
	}
}
