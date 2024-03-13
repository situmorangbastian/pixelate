.PHONY: run
run:
	go run app/main.go

.PHONY: test
test:
	GO111MODULE=on go test -covermode=atomic ./...
