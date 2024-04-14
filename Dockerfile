FROM golang:1.21.6

WORKDIR /app

COPY . .

RUN go mod download && go mod tidy

RUN go build -o main ./cmd/http/main.go

ENTRYPOINT ["./app/main"]