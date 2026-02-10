# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server


# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

COPY models ./models

EXPOSE 8080

CMD ["./server"]
