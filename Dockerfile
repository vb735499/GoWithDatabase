FROM golang:1.22

WORKDIR /app

COPY . /app

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]