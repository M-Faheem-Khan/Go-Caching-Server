PROJECT_NAME := cachingServer

build:
	go build -o ./bin/$(PROJECT_NAME) ./src/main.go

clean:
	rm -rf ./bin/
	mkdir bin

run:
	docker-compose up