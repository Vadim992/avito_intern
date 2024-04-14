FROM golang:1.21.6

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./cmd/http/main.go

ENTRYPOINT ["/app/main"]