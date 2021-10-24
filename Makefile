GO_SRCS := $(shell find . -type f -name '*.go')

main: queries/query.sql.go $(GO_SRCS)
	go build main.go

run: main
	./main

queries/query.sql.go: db/query.sql db/schema.sql
	sqlc generate

.PHONY: run
