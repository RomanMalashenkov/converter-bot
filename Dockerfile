FROM golang:latest

WORKDIR /go/src/tg_bot

COPY ./ ./

RUN go mod download
RUN go build -o tg_bot ./cmd/main.go

CMD ["./tg_bot"]