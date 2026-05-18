FROM golang:1.25.5-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o home_security main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/home_security .

EXPOSE 8080
CMD ["./home_security"]