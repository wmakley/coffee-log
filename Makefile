GO_SRCS := $(shell find . -type f -name '*.go')

main: migrate sqlc $(GO_SRCS)
	go build main.go

clean:
	rm -f ./main

run: main
	./main

sqlc:
	sqlc generate

test: migrate-test sqlc
	go test ./db/sqlc

migrate:
	dbmate up

migrate-test:
	dbmate -e TEST_DATABASE_URL up

rollback-test:
	dbmate -e TEST_DATABASE_URL down

clean-test:
	dbmate -e TEST_DATABASE_URL drop

.PHONY: run test sqlc migrate rollback migrate-test rollback-test drop-test
