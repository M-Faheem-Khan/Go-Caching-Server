FROM golang:latest

WORKDIR /cachingServer
COPY . .

RUN go build -o ./bin/cachingServer ./src/main.go

CMD ["./bin/cachingServer"]
