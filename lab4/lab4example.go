package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/net/html"

	"log"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func readItem(item *html.Node) *Item {
	return &Item{
		Ref:   getAttr(item, "href"),
		Title: getAttr(item, "title"),
	}
}

// Item is an item
type Item struct {
	Ref, Title string
}

func downloadNews(link string) []*Item {
	log.Printf("sending request to %s", link)
	if response, err := http.Get("http://" + link); err != nil {
		log.Panic(err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Printf("got response from %s : %d", link, status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Panic("invalid HTML from ", link, " error ", err)
			} else {
				log.Printf("HTML from %s parsed successfully", link)
				search(doc)
				return nil
			}
		}
	}
	return nil
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

func search(node *html.Node) []*Item {
	if isDiv(node, "title") {
		var items []*Item
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isElem(c, "a") {
				if getAttr(c, "data-tn-element") == "jobTitle" {
					println(getAttr(c, "title"))
					item := readItem(c)
					if item == nil {
						log.Println("appending, len = ", len(items))
						items = append(items, item)
					}
				}
			}
		}
		return items
	}
	var result []*Item
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, search(c)...)
	}
	return result
}

func main() {
	var vacancy string
	flag.StringVar(&vacancy, "vac", "менеджер", "wanted vacancy")
	flag.Parse()
	downloadNews(fmt.Sprintf("ru.indeed.com/%s-jobs", vacancy))
}
