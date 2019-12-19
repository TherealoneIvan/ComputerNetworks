package main

import (
	"fmt"
	"net/http"

	log "github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
)

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

// Item is a struct
type Item struct {
	Ref, Time, Title string
}

func isA(node *html.Node, class string) bool {
	return isElem(node, "ul") && getAttr(node, "class") == class
}

func readItem(item *html.Node, boolean bool) *Item {

	for c := item.FirstChild; c != nil; c = c.NextSibling {
		for _, i := range c.Attr {
			fmt.Println("http://lab-sud.ru" + i.Val[0:len(i.Val)-1])
			//fmt.Println("!")
			if boolean {
				downloadNews("http://lab-sud.ru"+i.Val[0:len(i.Val)], false)
			}
		}
	}

	return nil
}

func downloadNews(theme string, boolean bool) []*Item {
	log.Info("sending request to drive2.ru")
	if response, err := http.Get(theme); err != nil {
		fmt.Printf("error")
		log.Error("request to drive2.ru failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		//fmt.Printf("success1")
		log.Info("got response from drive2.ru", "status", status)
		if status == http.StatusOK {
			//fmt.Printf("success2")
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from drive2.ru", "error", err)
			} else {
				log.Info("HTML from drive2.ru parsed successfully")
				//fmt.Println(theme)

				return search(doc, boolean)
			}
		}
	}
	return nil
}

func search(node *html.Node, boolean bool) []*Item {
	kek := isLi(node.Parent, "next")
	var all = make([]*Item, 0)
	if boolean && isA(node, "rubrics") {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			all = append(all, readItem(c, true))
		}
	} else if !boolean && kek && isA(node, "rubrics") {
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			all = append(all, readItem(c, false))
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c, boolean); items != nil {
			for _, k := range items {
				all = append(all, k)
			}
		}
	}

	return all
}

func isLi(node *html.Node, class string) bool {
	return isElem(node, "li") && getAttr(node, "class") == class
}

//===================================================================================================

func main() {

	log.Info("Downloader started")
	var items []*Item
	items = downloadNews("http://lab-sud.ru", true)

	fmt.Printf("%d ", len(items))
	//for _, i := range items {
	//mt.Printf("%s %s %s\n", i.Ref, i.Time, i.Title)
	//sitems.
	//}
}
