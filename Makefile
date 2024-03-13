.PHONY: run
run:
	go run app/main.go

.PHONY: test
test:
	GO111MODULE=on go test -covermode=atomic ./...

.PHONY: docker-build
docker-build:
	docker build -t pixelate .

.PHONY: docker-run
docker-run: docker-build
	docker run -it --rm -p 1111:1111 pixelate
