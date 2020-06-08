FROM golang:1.14-alpine

RUN apk update && apk add git && apk add bash

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
