package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Link struct {
	Text    string // текст ссылки
	URL     string // внешняя ссылка
	PageURL string // URL страницы сайта, где находится внешняя ссылка
}

var externalLinks = []*Link{}

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains:
		colly.AllowedDomains("hazadus.ru"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Ссылки, которые начинаются не с https://hazadus.ru, https://amgold.ru, mailto:,
		// добавлять в массив
		if isExternalLink(link) {
			fmt.Printf("    Найдена внешняя ссылка: %q -> %s\n", e.Text, link)
			addExternalLink(&Link{
				Text:    e.Text,
				URL:     link,
				PageURL: e.Request.URL.String(),
			})
		}

		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	c.Visit("https://hazadus.ru/")

	// Внешние ссылки в массиве обойти, для каждой зафиксировать успех/неуспех
	// Неуспешные добавить в массив broken
	var brokenLinks = []*Link{}
	for i, l := range externalLinks {
		fmt.Printf("\n%d. %q - %s\n", i+1, l.Text, l.URL)

		if !checkStatusAndRedirects(l.URL) {
			fmt.Printf("    Ошибка при обращении к URL\n    @ %q\n", l.PageURL)
			brokenLinks = append(brokenLinks, l)
		}
	}

	fmt.Println("БИТЫЕ ССЫЛКИ")
	for i, l := range brokenLinks {
		fmt.Printf("%d. %q - %s - %q\n", i+1, l.PageURL, l.Text, l.URL)
		// TODO: формировать Markdown-файл с отчётом в файл (скриптом GHA потом его переименуем)
	}
}

func isExternalLink(url string) bool {
	return strings.HasPrefix(url, "http") &&
		(!strings.HasPrefix(url, "https://hazadus.ru") || !strings.HasPrefix(url, "https://amgold.ru") || !strings.HasPrefix(url, "mailto:"))
}

func addExternalLink(link *Link) {
	for _, l := range externalLinks {
		if l.URL == link.URL {
			return
		}
	}

	externalLinks = append(externalLinks, link)
}

func checkStatusAndRedirects(url string) bool {
	var redirectCount int
	nextURL := url
	maxRedirectsAllowed := 100
	getRequestTimeout := 3 * time.Second

	// Считаем количество редиректов
	for redirectCount <= maxRedirectsAllowed {
		httpClient := http.Client{
			Timeout: getRequestTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		response, err := httpClient.Get(nextURL)
		if err != nil {
			// TODO: возвращать описание ошибки
			fmt.Println(err)
			return false
		}
		if response.StatusCode == 200 {
			return true
		} else {
			nextURL = response.Header.Get("Location")
			redirectCount += 1
		}
	}

	return false
}
