alias t := test
alias ti := test-integration
alias tu := test-unit
alias l := lint
alias run := run-app

# Lists all available tasks
default:
    @just --list

dockerbuild:
    docker --debug build -t cqrs-monitored-app -f build/Dockerfile .

build:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"	\
    -o bin/cqrs-monitored-app \
    cmd/app/main.go

# Run the app with hot reload
air:
    air -c app.air.toml

# Runs the application with the .env.example file
run-app:
    cp .env.example .env
    go run cmd/app/main.go

# Run all tests, or any tests specified by the path with its extra parameters
test path="./..." *params="":
    go test {{path}} -race {{params}}

# Runs all tests located at ./test
test-integration *params:
    @just test ./test/... -timeout 300s {{params}}

# Runs all tests, except integration tests located at ./test
test-unit *params:
    go test -p 2 $(go list ./... | grep -v ./test) -race {{params}}

# Clears the test cache
clear-cache:
    go clean -testcache

# Formats all go files
lint:
    go fmt ./...