package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/RealJK/rss-parser-go"
)

//VacancyInfo represents data of a vacancy
type VacancyInfo struct {
	Title     string `json:"title"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
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

func getVacancyInfo(divjson []byte) []byte {
	// десереализация JSON ответа от сервера сайта "Мой круг"
	var wrapper DivisionWrapper
	json.Unmarshal(divjson, &wrapper)

	// парсинг RSS-ленты вакансий и сериализация результата в JSON
	var vacancies []VacancyInfo
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
						Title:     v.Title[strings.Index(v.Title, "«"):strings.Index(v.Title, "»")],
						Timestamp: string(v.PubDate),
						Type:      div.Alias,
					})
				}
			}
		}
	}
	vinfojson, _ := json.Marshal(vacancies)

	return vinfojson
}

// VInfoRouterHandler : Отдаёт ленту вакансий с сайта moikrug в формате json
func VInfoRouterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(getVacancyInfo(getDivisions()))
}

// DivsRouterHandler : Отдаёт список сфер деятельности с сайта moikrug в формате json
func DivsRouterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// десереализация JSON ответа от сервера сайта "Мой круг"
	divjson := getDivisions()
	var wrapper DivisionWrapper
	json.Unmarshal(divjson, &wrapper)
	result, _ := json.Marshal(wrapper.Content)
	w.Write(result)
}

func main() {
	http.HandleFunc("/getvinfo", VInfoRouterHandler)
	http.HandleFunc("/getdivisions", DivsRouterHandler)
	err := http.ListenAndServe(":8017", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
