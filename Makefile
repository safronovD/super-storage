include variables.mk

.PHONY: all

all: version dependency lint build docker-build docker-run

# print version
version:
	@printf $(TAG)

dependency:
	${GO_ENV_VARS} go mod download

### Lint

lint:
	${GO_ENV_VARS} golangci-lint -v run ./...

### Unit tests

coverage:
	go tool cover -html=coverage.out -o coverage.html

test:
	${GO_ENV_VARS} go test `go list ./... | grep pkg` -race -cover -coverprofile=coverage.out -covermode=atomic

### Build

clean:
	rm -rf ./build

build: clean
	${GO_ENV_VARS} go build -a -o ./build/storage ./cmd/main.go

### Docker

docker-build:
	docker build --tag ${IMAGE_NAME}:${TAG} .

docker-run:
	docker run ${IMAGE_NAME}:${TAG}
