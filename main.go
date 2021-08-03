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

func getSummary(url string) string {
	resp, err := soup.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	doc := soup.HTMLParse(resp)
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
	resp, err := soup.Get(baseUrl + "/news")
	if err != nil {
		log.Fatal(err)
	}
	doc := soup.HTMLParse(resp)
	for _, link := range doc.Find("div", "class", "nw-c-most-read").FindAll("a", "class", "gs-c-promo-heading") {
		title := link.FullText()
		url := baseUrl + link.Attrs()["href"]
		article := Article{
			Title:   title,
			URL:     url,
			Summary: getSummary(url),
		}
		articles = append(articles, article)
	}
	return articles
}

func main() {
	for _, article := range getArticles() {
		fmt.Println("### [" + article.Title + "](" + article.URL + ")\n" + article.Summary)
	}
}
