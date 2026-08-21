# Stage 1: build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/main ./cmd/api

# Stage 2: imagem final
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main

RUN chmod +x /app/main

EXPOSE 8082

CMD ["/app/main"]