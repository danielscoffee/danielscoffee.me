GO_BIN := $(shell go env GOPATH)/bin
TEMPL := $(shell command -v templ 2>/dev/null || echo $(GO_BIN)/templ)
AIR := $(shell command -v air 2>/dev/null || echo $(GO_BIN)/air)
TAILWIND_SHA256 := 7d24f7fa191d2193b78cd5f5a42a6093e14409521908529f42d80b11fde1f1d4

all: build test

templ-install:
	@if ! command -v templ > /dev/null && [ ! -x "$(GO_BIN)/templ" ]; then \
		echo "templ not found. Installing..."; \
		go install github.com/a-h/templ/cmd/templ@v0.3.1020; \
	fi

tailwind-install:
	@if [ ! -f tailwindcss ]; then \
		trap 'rm -f tailwindcss.tmp' EXIT; \
		curl -fL https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64 -o tailwindcss.tmp; \
		printf '%s  %s\n' "$(TAILWIND_SHA256)" tailwindcss.tmp | sha256sum -c -; \
		mv tailwindcss.tmp tailwindcss; \
	fi
	@printf '%s  %s\n' "$(TAILWIND_SHA256)" tailwindcss | sha256sum -c -
	@chmod +x tailwindcss

generate: templ-install tailwind-install
	@$(TEMPL) generate -path .
	@./tailwindcss -i internal/web/styles/input.css -o internal/web/assets/css/output.css

build: generate
	@go build -o main cmd/api/main.go

run: generate
	@go run cmd/api/main.go

test: generate
	@go test ./... -v

docker-run:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose up --build; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

docker-down:
	@if docker compose version >/dev/null 2>&1; then \
		docker compose down; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

clean:
	@rm -f main

watch:
	@if ! command -v air > /dev/null && [ ! -x "$(GO_BIN)/air" ]; then \
		echo "air not found. Installing..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@echo "Watching..."
	@$(AIR)

.PHONY: all build run test clean watch tailwind-install docker-run docker-down templ-install generate
