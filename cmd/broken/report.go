package main

import (
	"fmt"
	"os"
)

// createMarkdownReport генерирует отчёт о битых ссылках из списка outputFilePath в формате Markdown,
// и сохраняет его в файл outputFilePath.
func createMarkdownReport(brokenLinks []*Link, outputFilePath string) error {
	report := fmt.Sprintf("# ⛓️‍💥 Битые ссылки\n\nВсего: %d\n", len(brokenLinks))

	for i, l := range brokenLinks {
		report += fmt.Sprintf("\n## %d. На странице %s\n\n- ⌨️ Текст ссылки: %q\n- ⛓️‍💥 Внешний URL: %s\n- ⚠️ Ошибка: %s\n",
			i+1,
			l.PageURL,
			l.Text,
			l.URL,
			l.Error,
		)
	}

	return os.WriteFile(outputFilePath, []byte(report), 0444)
}
