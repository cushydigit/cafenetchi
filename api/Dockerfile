FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod download

# BUILD the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/api ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/api /app/api

EXPOSE 8080

CMD ["/app/api"]

