# Build stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app
COPY . .

RUN go build -o main main.go

# Run Stage

FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080
CMD ["/app/main"]
