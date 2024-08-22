#steps to run the application

go mod tidy
go build ./cmd/server
go run ./cmd/server/main.go