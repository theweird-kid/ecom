build:
	@go build -o bin/ecom cmd/main.go
test:
	@go test -v ./...
run: build
	@./bin/ecom
migrate:
	@cd ./cmd/migrations/schema && goose postgres postgres://postgres:postgres@localhost:5432/ecom_db up