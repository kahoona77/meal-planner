FROM golang:1.27-alpine AS builder
# keep in sync with mise.toml
ARG TAILWIND_VERSION=4.3.3
RUN apk add --no-cache gcc musl-dev
ADD --chmod=755 https://github.com/tailwindlabs/tailwindcss/releases/download/v${TAILWIND_VERSION}/tailwindcss-linux-x64-musl /usr/local/bin/tailwindcss
WORKDIR /go/src/app
COPY . .
RUN tailwindcss -i web/styles/main.css -o web/static/main.css --minify
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -mod=vendor -ldflags '-linkmode external -extldflags "-static"' -o meal-planner
RUN go test -mod=vendor -v ./...

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY web/tmpl/ /app/web/tmpl/
COPY --from=builder /go/src/app/web/static/ /app/web/static/
COPY --from=builder /go/src/app/meal-planner /app/

ENTRYPOINT ["/app/meal-planner"]

EXPOSE 8080
