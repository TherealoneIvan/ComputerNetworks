# Лабораторная работа № 0.1

## Разработка простейшего web-сервера

Рассматривается задача разработки web-сервера на языке GO на основе пакета net/http.

### Пример реализации простого web-сервера

```go
package main

import (
	"fmt"      // пакет для форматированного ввода вывода
	"log"      // пакет для логирования
	"net/http" // пакет для поддержки HTTP протокола
	"strings"  // пакет для работы с  UTF-8 строками
)

// HomeRouterHandler : обработчик домашней страницы
func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //анализ аргументов,
	fmt.Println(r.Form) // ввод информации о форме на стороне сервера
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Test!") // отправляем данные на клиентскую сторону
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)  // установим роутер
	err := http.ListenAndServe(":9000", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```

**Задача 1:** Реализовать web-сервер и запустить на заданном порте.

**Задача 2:** Изучить принимаемые web-сервером параметры, реализовать передачу данных методом GET.

**Задача 3:** Реализовать вывод форматированного гипертекста по вариантам:
- контекстное меню в виде гиперссылок, при клике на гиперссылку должна выполняться подмена контента;
- вывод формы;
- вывод фреймового окна.

**Замечание 1:** Для корректной работы GO на серверах 185.20.227.83, 185.20.226.174 необходимо задать переменную окружения 
```bash
export GOPATH=~/go.
```

**Замечание 2:** Запуск приложение в режиме интерпретатора выполняется следующим образом
```bash
go run example_app.go.
```

# Лабораторная работа № 0.2

## Разработка приложения обработки данных из RSS-канала

Рассматривается задача разработки приложения на языке GO реализующего синтаксический разбор XML файла формата RSS.

Для реализации данной задачи предлагается использовать библиотеку rss-parser-go, которая доступна по адресу https://github.com/masterjk/rss-parser-go.

Установка библиотеки  rss-parser-go:
```bash
go get github.com/RealJK/rss-parser-go
```

### Пример использования библиотеки

```go
package main

import (
	"fmt"

	"github.com/RealJK/rss-parser-go"
)

func main() {

	rssObject, err := rss.ParseRSS("http://blagnews.ru/rss_vk.xml")
	if err != nil {

		fmt.Printf("Title           : %s\n", rssObject.Channel.Title)
		fmt.Printf("Generator       : %s\n", rssObject.Channel.Generator)
		fmt.Printf("PubDate         : %s\n", rssObject.Channel.PubDate)
		fmt.Printf("LastBuildDate   : %s\n", rssObject.Channel.LastBuildDate)
		fmt.Printf("Description     : %s\n", rssObject.Channel.Description)

		fmt.Printf("Number of Items : %d\n", len(rssObject.Channel.Items))

		for v := range rssObject.Channel.Items {
			item := rssObject.Channel.Items[v]
			fmt.Println()
			fmt.Printf("Item Number : %d\n", v)
			fmt.Printf("Title       : %s\n", item.Title)
			fmt.Printf("Link        : %s\n", item.Link)
			fmt.Printf("Description : %s\n", item.Description)
			fmt.Printf("Guid        : %s\n", item.Guid.Value)
		}
	}
}
```

**Замечание 1:** Для корректной работы GO на серверах 185.20.227.83, 185.20.226.174 необходимо задать переменную окружения 
```bash
export GOPATH=~/go.
```

**Замечание 2:** Запуск приложение в режиме интерпретатора выполняется следующим образом
```bash
go run example_app.go.
```

**Задача:** Реализовать получение данных из различных RSS- каналов по вариантам. Сравнить результаты разбора и сделать выводы.

- http://blagnews.ru/rss_vk.xml

- http://www.rssboard.org/files/sample-rss-2.xml

- https://lenta.ru/rss

- https://news.mail.ru/rss/90/

- http://technolog.edu.ru/index.php?option=com_k2&view=itemlist&layout=category&task=category&id=8&lang=ru&format=feed

- https://vz.ru/rss.xml

- http://news.ap-pa.ru/rss.xml

# Лабораторная работа № 0.3

## Разработка web-ориентированного клиент-серверного приложения получения и представления данных из RSS-канала

Целью данной лабораторной работы является произвести интеграцию результатов работ проведенных в лабораторной работе № 0.1 и лабораторной работе № 0.2.

**Задача:** необходимо разработать web-сервер, который выполняет соединение с удаленным (удаленными) серверами RSS-новостей и возвращает результаты обработки данных в структурированном виде (страница гипертекста) web-клиенту, в нашем случае в браузер по вариантам.

**Замечание 1:** Для корректной работы GO на серверах 185.20.227.83, 185.20.226.174 необходимо задать переменную окружения 
```bash
export GOPATH=~/go.
```

**Замечание 2:** Запуск приложение в режиме интерпретатора выполняется следующим образом
```bash
go run example_app.go.
```