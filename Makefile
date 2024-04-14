run: ./cmd/http/main.go
	go run ./cmd/http/main.go

test: ./internal/requests_test.go
	go test -v ./...
up:
	docker compose up