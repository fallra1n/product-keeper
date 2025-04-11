FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /app/product-keeper ./cmd/product-keeper/main.go

FROM alpine:latest
WORKDIR /app/

ARG JWT_SECRET
ARG CONFIG_PATH
ENV JWT_SECRET=${JWT_SECRET}
ENV CONFIG_PATH=${CONFIG_PATH}

COPY --from=builder /app/product-keeper ./
COPY --from=builder /app/config ./config/
COPY --from=builder /app/.cert ./.cert/

CMD ["./product-keeper"]