package main

import (
	"fmt"
	"log"

	"github.com/anaskhan96/soup"
)

type Article struct {
	URL     string
	Title   string
	Summary string
}

func getHtml(url string) soup.Root {
	resp, err := soup.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return soup.HTMLParse(resp)
}

func getSummary(url string) string {
	doc := getHtml(url)
	link := doc.Find("p", "role", "introduction")
	if link.Error == nil {
		return link.FullText()
	}
	for _, link := range doc.Find("article").Children() {
		if link.Attrs()["data-component"] == "text-block" {
			return link.FullText()
		}
	}
	return ""
}

func getArticles() []Article {
	var articles []Article
	baseUrl := "https://www.bbc.com"
	doc := getHtml(baseUrl + "/news")
	for _, link := range doc.Find("div", "class", "nw-c-most-read").FindAll("a", "class", "gs-c-promo-heading") {
		url := baseUrl + link.Attrs()["href"]
		articles = append(articles, Article{
			Title:   link.FullText(),
			URL:     url,
			Summary: getSummary(url),
		})
	}
	return articles
}

func main() {
	for _, article := range getArticles() {
		fmt.Println("### [" + article.Title + "](" + article.URL + ")\n" + article.Summary)
	}
}
