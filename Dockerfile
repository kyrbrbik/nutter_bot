FROM golang:1.20.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o main .

FROM cgr.dev/chainguard/wolfi-base:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
CMD ["/app/main"]

