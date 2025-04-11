FROM golang:1.24-alpine

WORKDIR /app

ARG JWT_SECRET
ARG CONFIG_PATH

ENV JWT_SECRET=${JWT_SECRET}
ENV CONFIG_PATH=${CONFIG_PATH}

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o /app/product-keeper ./cmd/product-keeper/main.go

CMD ["/app/product-keeper"]
