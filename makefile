run:
	@go run cmd/main.go

build:
	@go build -o bin/cachecast cmd/main.go

docker-run:
	@docker compose up  

docker-build:
	@docker compose build 

