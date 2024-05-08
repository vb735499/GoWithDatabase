FROM golang:1.22

WORKDIR /app

COPY . /app

RUN go mod init github.com/my/repo
RUN go get github.com/redis/go-redis/v9

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
