.PHONY: build run test clean docker-up docker-down

build:
	go build -o main .

run: build
	./main

test:
	go test ./...

clean:
	rm -f main

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

test-api:
	chmod +x test_api.sh
	./test_api.sh

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...
