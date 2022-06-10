FROM golang:1.18

WORKDIR /usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/app ./cmd/app/main.go

EXPOSE 8000

CMD ["app"]
