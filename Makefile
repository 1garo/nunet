svr:
	@go build -o svr ./server

cli:
	@go build -o cli ./client/main.go

deps:
	@go mod tidy

clean:
	@go clean -i

up:
	@docker compose build
	@docker compose up -d

down:
	@docker compose down --remove-orphans

test:
	@go test ./...

testv:
	@go test -v ./...


