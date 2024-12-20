## help: вывести это сообщение
.PHONY: help
help:
	@echo 'Команды для make:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## format: отформатировать весь код Go в поддиректориях
.PHONY: format
format:
	go fmt ./...

## lint: запустить golangci-lint (должен быть установлен)
.PHONY: lint
lint:
	golangci-lint run -c golangci.yml

## test: запустить все тесты
.PHONY: test
test:
	go test ./...

## test/cover: запустить тесты и показать покрытие
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

## build: скомпилировать исполняемый файл
.PHONY: build
build:
	go build -ldflags=${linker_flags} -o ./bin/broken ./cmd/broken/

## run: скомпилировать и запустить исполняемый файл
.PHONY: run
run: build
	./bin/broken


## todo: вывести перечень всех комментов с меткой TODO в файлах проекта
.PHONY: todo
todo:
	grep -r "TODO:" .

## cloc: посчитать строки кода в проекте и сохранить в файл
.PHONY: cloc
cloc:
	cloc . > cloc.txt
