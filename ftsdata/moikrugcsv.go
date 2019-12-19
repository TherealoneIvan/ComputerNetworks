package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/RealJK/rss-parser-go"
	"github.com/gocarina/gocsv"
)

//VacancyInfo represents data of a vacancy
type VacancyInfo struct {
	ID          int    `csv:"id"`
	Title       string `csv:"title"`
	Timestamp   string `csv:"timestamp"`
	Type        string `csv:"type"`
	Description string `csv:"description"`
}

//Division is a filed of vacancy
type Division struct {
	Name  string `json:"name"`
	Alias string `json:"alias_name"`
}

// DivisionWrapper wraps a slice of divisions
type DivisionWrapper struct {
	Content []Division `json:"divisions"`
}

func getDivisions() []byte {
	// сделать GET запрос к API сайта "Мой круг"
	response, err := http.Get("https://api.moikrug.ru/v1/integrations/divisions?access_token=01728013a28b7a64c9be6c59b6d76757ff426fd8b0aef9711fdfd2280ad66516")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// занести в массив байтов тело ответа от сервера
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func getVacancyInfo(divjson []byte) []VacancyInfo {
	// десереализация JSON ответа от сервера сайта "Мой круг"
	var wrapper DivisionWrapper
	json.Unmarshal(divjson, &wrapper)

	// парсинг RSS-ленты вакансий и сериализация результата в JSON
	var vacancies []VacancyInfo
	j := 1
	for _, div := range wrapper.Content {
		i := 1
		for {
			rssObject, err := rss.ParseRSS(fmt.Sprintf("https://moikrug.ru/vacancies/rss?divisions=%s&page=%d", div.Alias, i))
			i++
			if len(rssObject.Channel.Items) == 0 {
				break
			}
			if err != nil {
				for _, v := range rssObject.Channel.Items {
					vacancies = append(vacancies, VacancyInfo{
						ID:          j,
						Title:       v.Title[strings.Index(v.Title, "«")+2 : strings.Index(v.Title, "»")],
						Timestamp:   string(v.PubDate),
						Type:        div.Alias,
						Description: v.Description,
					})
					j++
				}
			}
		}
	}

	return vacancies
}

func main() {
	// получаем список всех вакансий
	vacancies := getVacancyInfo(getDivisions())

	// создадим выходной файл
	vacanciesFile, err := os.OpenFile("vacancies.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer vacanciesFile.Close()

	// настроим сериализацию в CSV
	gocsv.SetCSVWriter(func(out io.Writer) *gocsv.SafeCSVWriter {
		writer := csv.NewWriter(out)
		writer.Comma = '^'
		return gocsv.NewSafeCSVWriter(writer)
	})

	// запись в файл
	err = gocsv.MarshalFile(&vacancies, vacanciesFile)
	if err != nil {
		panic(err)
	}
}
