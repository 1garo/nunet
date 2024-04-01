build-s:
	@go build -o nu-server ./server/main.go

build-c:
	@go build -o nu-client ./client/main.go

deps:
	@go mod tidy

clean:
	@go clean -i

up:
	@docker compose build
	@docker compose up -d

down:
	@docker compose down --remove-orphans


