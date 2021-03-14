git:
	git add .
	git commit -m "$(msg)"
	git push origin master

fmt:
	@echo "formatting code"
	go fmt ./...

lint:
	@echo "Linting source code"
	golint ./...

vet:
	@echo "Checking for code issues"
	go vet ./...

test:
	@echo "running tests"
	go test ./...

install:
	@echo "installing external dependencies"
	go mod download

graphql:
	@echo "generating graphql stubs"
	go run github.com/99designs/gqlgen generate

run:
	 go run cmd/main.go

dev:
	docker-compose -f dev.yml up -d