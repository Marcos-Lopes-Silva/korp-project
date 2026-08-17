# services/stock/Dockerfile

# Stage 1: build
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: imagem final, só o binário
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8082
CMD ["./main"]