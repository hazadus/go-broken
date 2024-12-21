package main

import (
	"errors"
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
	Error   string // Текст ошибки
}

// collectExternalLinks собирает все "внешние" ссылки с указанного сайта.
//   - siteUrl - страница, с которой начинаем собирать ссылки
//   - allowedDomain - домен, с которого разрешается собирать ссылки
//   - internalURLs - список URL, которые считаются "внутренними"
func collectExternalLinks(siteUrl, allowedDomain string, internalURLs []string) ([]*Link, error) {
	var externalLinks = []*Link{}

	c := colly.NewCollector(
		// Visit only domains:
		colly.AllowedDomains(allowedDomain),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Ссылки, которые начинаются не с "внутренних", добавлять в массив
		if isExternalLink(link, internalURLs) {
			fmt.Printf("    Найдена внешняя ссылка: %q -> %s\n", e.Text, link)
			if !isIncluded(externalLinks, link) {
				externalLinks = append(externalLinks, &Link{
					Text:    e.Text,
					URL:     link,
					PageURL: e.Request.URL.String(),
				})
			}
		}

		// Visit link found on page. Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Начать сбор ссылок
	c.Visit(siteUrl)

	return externalLinks, nil
}

// checkLinks проверяет каждую ссылку из externalLinks при помощи checkStatusAndRedirects.
// "Битые" ссылки возвращаются в массиве.
func checkLinks(externalLinks []*Link) ([]*Link, error) {
	var brokenLinks = []*Link{}
	for i, l := range externalLinks {
		fmt.Printf("%d. %q - %s\n", i+1, l.Text, l.URL)

		ok, err := checkStatusAndRedirects(l.URL)
		if !ok {
			if err != nil {
				fmt.Printf("    Ошибка: %s\n    @ %q\n", err.Error(), l.PageURL)
				l.Error = err.Error()
			} else {
				fmt.Printf("    Ошибка соединения с URL\n    @ %q\n", l.PageURL)
			}

			brokenLinks = append(brokenLinks, l)
		}
	}

	return brokenLinks, nil
}

// isExternalLink проверяет, является ли ссылка "внешней" для сканируемого сайта.
func isExternalLink(url string, interalURLs []string) bool {
	if !strings.HasPrefix(url, "http") {
		return false
	}

	for _, internalURL := range interalURLs {
		if strings.HasPrefix(url, internalURL) {
			return false
		}
	}

	return true
}

// isIncluded проверяет наличие ссылки в массиве.
func isIncluded(externalLinks []*Link, url string) bool {
	for _, l := range externalLinks {
		if l.URL == url {
			return true
		}
	}

	return false
}

// checkStatusAndRedirects проверяет доступность ссылки с учетом редиректов.
// Будет происходить переход по редиректам, пока не получим статус-код 200.
// Возвращает true в случае успеха, false если статус-кода 200 так и не дождались.
func checkStatusAndRedirects(url string) (bool, error) {
	var redirectCount int
	nextURL := url
	maxRedirectsAllowed := 100
	getRequestTimeout := 5 * time.Second

	for redirectCount <= maxRedirectsAllowed {
		httpClient := http.Client{
			Timeout: getRequestTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}

		response, err := httpClient.Get(nextURL)
		if err != nil {
			return false, err
		}

		if response.StatusCode == 200 {
			return true, nil
		} else if response.StatusCode >= 400 {
			return false, errors.New(fmt.Sprintf("GET %q → status code = %d", nextURL, response.StatusCode))
		} else {
			url, err := response.Location()
			if err != nil {
				return false, err
			}

			nextURL = url.String()
			redirectCount += 1
		}
	}

	return false, nil
}
