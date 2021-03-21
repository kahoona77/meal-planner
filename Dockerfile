FROM golang:1.16-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /go/src/app
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -mod=vendor -ldflags '-linkmode external -extldflags "-static"' -o meal-planner
RUN go test -mod=vendor -v ./...


FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app
# copy binary into image
COPY web/tmpl/ /app/web/tmpl/
COPY web/assets/ /app/web/assets/
COPY --from=builder /go/src/app/meal-planner /app/

ENTRYPOINT ["/app/meal-planner"]

EXPOSE 8080