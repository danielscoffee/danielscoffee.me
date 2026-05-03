all: build test

define docker_compose
	@if docker compose $(1) 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose $(1); \
	fi
endef

templ-install:
	@if ! command -v templ > /dev/null; then \
		echo "templ not found. Installing..."; \
		go install github.com/a-h/templ/cmd/templ@latest; \
	fi

tailwind-install:
	@if [ ! -f tailwindcss ]; then \
		curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tailwindcss; \
	fi
	@chmod +x tailwindcss

generate: templ-install tailwind-install
	@templ generate -path .
	@./tailwindcss -i internal/web/styles/input.css -o internal/web/assets/css/output.css

build: generate
	@go build -o main cmd/api/main.go

run: generate
	@go run cmd/api/main.go

test: generate
	@go test ./... -v

docker-run:
	$(call docker_compose,up --build)

docker-down:
	$(call docker_compose,down)

clean:
	@rm -f main

watch:
	@if command -v air > /dev/null; then \
		echo "Watching..."; \
		air; \
	else \
		echo "air not found. Installing..."; \
		go install github.com/air-verse/air@latest; \
		echo "Watching..."; \
		air; \
	fi

.PHONY: all build run test clean watch tailwind-install docker-run docker-down templ-install generate
