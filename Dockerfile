FROM golang:alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

ADD . .

CMD ["go","run","cmd/app/main.go"]
