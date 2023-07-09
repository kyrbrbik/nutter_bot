FROM golang:1.20.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -o main .

FROM debian:11-slim
RUN apt update && apt install ca-certificates -y
WORKDIR /app
COPY --from=builder /app/main .
ARG DISCORD_TOKEN
ARG OPENAI_TOKEN
CMD ["/app/main"]

