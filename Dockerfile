FROM golang:1.22-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o server cmd/main.go
FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server /app/
COPY config/config.yaml /app/config/config.yaml
EXPOSE 8080
ENTRYPOINT ["/app/server"]
