package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	siteURLFlag := flag.String("s", "https://hazadus.ru", "страница, с которой начинаем собирать ссылки")
	allowedDomainFlag := flag.String("d", "hazadus.ru", "домен, с которого разрешается собирать ссылки")
	internalURLsFlag := flag.String("i", "https://hazadus.ru,https://amgold.ru", "список URL, которые считаются внутренними, через запятую")
	flag.Parse()

	internalURLs := strings.Split(*internalURLsFlag, ",")

	run(*siteURLFlag, *allowedDomainFlag, internalURLs)
}

// run выполняет основную логику приложения.
func run(siteURL, allowedDomain string, internalURLs []string) {
	externalLinks, err := collectExternalLinks(
		siteURL,
		allowedDomain,
		internalURLs,
	)
	if err != nil {
		fmt.Printf("Ошибка при сборе внешних ссылок: %s\n", err)
		os.Exit(1)
	}

	brokenLinks, err := checkLinks(externalLinks)
	if err != nil {
		fmt.Printf("Ошибка при проверке внешних ссылок: %s\n", err)
		os.Exit(1)
	}

	err = createMarkdownReport(brokenLinks, "./report.md")
	if err != nil {
		fmt.Printf("Ошибка при сохранении отчета: %s\n", err)
		os.Exit(1)
	}
}
