FROM golang:1.26.5-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS build
RUN apk add --no-cache curl libstdc++ libgcc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1020 && \
    templ generate -path . && \
    curl -fL https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64 -o tailwindcss && \
    echo "7d24f7fa191d2193b78cd5f5a42a6093e14409521908529f42d80b11fde1f1d4  tailwindcss" | sha256sum -c - && \
    chmod +x tailwindcss && \
    ./tailwindcss -i internal/web/styles/input.css -o internal/web/assets/css/output.css

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main cmd/api/main.go

FROM gcr.io/distroless/static-debian12:nonroot@sha256:b7bb25d9f7c31d2bdd1982feb4dafcaf137703c7075dbe2febb41c24212b946f AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/content /app/content
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/main"]


