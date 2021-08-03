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
	r, err := soup.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	d := soup.HTMLParse(r)
	link := d.Find("p", "role", "introduction")
	if link.Error == nil {
		return link.FullText()
	}
	for _, link := range d.Find("article").Children() {
		dataComponent := link.Attrs()["data-component"]
		if dataComponent == "text-block" {
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
		if err != nil {
			log.Fatal(err)
		}
		article := Article{
			Title:   title,
			URL:     url,
			Summary: getSummary(url),
		}
		articles = append(articles, article)
	}
	return articles
}

func getMarkdown(articles []Article) string {
	s := ""
	for _, article := range articles {
		s += "### [" + article.Title + "](" + article.URL + ")\n" + article.Summary + "\n"
	}
	return s
}

func main() {
	articles := getArticles()
	fmt.Println(getMarkdown(articles))
}
